package planner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
)

// LLMSettings configures the optional content-adaptation LLM. The LLM never
// decides engagement state or intervention; it only rewrites activity content.
type LLMSettings struct {
	Provider string // "anthropic", "openai" (any OpenAI-compatible endpoint), or "gemini"
	APIKey   string
	Model    string
	BaseURL  string
}

// LLMStatus is the safe-to-expose view of the configuration (never the key).
type LLMStatus struct {
	Configured bool   `json:"configured"`
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	BaseURL    string `json:"base_url"`
	LastSource string `json:"last_source"` // "ai" or "rules"
	LastError  string `json:"last_error,omitempty"`
}

var (
	llmMu         sync.RWMutex
	llmRuntime    *LLMSettings // nil => fall back to environment
	llmLastSource = "rules"
	llmLastError  string
)

func envSettings() LLMSettings {
	s := LLMSettings{
		Provider: os.Getenv("LLM_PROVIDER"),
		APIKey:   os.Getenv("LLM_API_KEY"),
		Model:    os.Getenv("LLM_MODEL"),
		BaseURL:  os.Getenv("LLM_BASE_URL"),
	}
	if s.APIKey == "" {
		if k := os.Getenv("ANTHROPIC_API_KEY"); k != "" {
			s.APIKey = k
			if s.Provider == "" {
				s.Provider = "anthropic"
			}
		} else if k := os.Getenv("OPENAI_API_KEY"); k != "" {
			s.APIKey = k
		} else if k := os.Getenv("GEMINI_API_KEY"); k != "" {
			s.APIKey = k
			if s.Provider == "" {
				s.Provider = "gemini"
			}
		}
	}
	return s
}

func resolveDefaults(s LLMSettings) LLMSettings {
	s.Provider = strings.ToLower(strings.TrimSpace(s.Provider))
	if s.Provider != "anthropic" && s.Provider != "gemini" {
		s.Provider = "openai"
	}
	s.BaseURL = strings.TrimRight(strings.TrimSpace(s.BaseURL), "/")
	s.Model = strings.TrimSpace(s.Model)
	s.APIKey = strings.TrimSpace(s.APIKey)
	if s.Provider == "anthropic" {
		if s.BaseURL == "" {
			s.BaseURL = "https://api.anthropic.com"
		}
		if s.Model == "" {
			s.Model = "claude-sonnet-5-5"
		}
	} else if s.Provider == "gemini" {
		if s.BaseURL == "" {
			s.BaseURL = "https://generativelanguage.googleapis.com"
		}
		if s.Model == "" {
			s.Model = "gemini-2.0-flash"
		}
	} else {
		if s.BaseURL == "" {
			s.BaseURL = "https://api.openai.com/v1"
		}
		if s.Model == "" {
			s.Model = "gpt-4o-mini"
		}
	}
	return s
}

// CurrentLLMSettings returns runtime settings if set, otherwise environment ones.
func CurrentLLMSettings() LLMSettings {
	llmMu.RLock()
	defer llmMu.RUnlock()
	if llmRuntime != nil {
		return resolveDefaults(*llmRuntime)
	}
	return resolveDefaults(envSettings())
}

// SetLLMSettings stores settings in memory only. A blank key keeps the current one.
func SetLLMSettings(s LLMSettings) {
	cur := CurrentLLMSettings()
	llmMu.Lock()
	defer llmMu.Unlock()
	if strings.TrimSpace(s.APIKey) == "" {
		s.APIKey = cur.APIKey
	}
	llmRuntime = &s
}

// ClearLLMSettings disables the LLM (ignores environment keys too).
func ClearLLMSettings() {
	llmMu.Lock()
	defer llmMu.Unlock()
	llmRuntime = &LLMSettings{}
}

func RecordLLMOutcome(source, errMsg string) {
	llmMu.Lock()
	defer llmMu.Unlock()
	llmLastSource, llmLastError = source, errMsg
}

func GetLLMStatus() LLMStatus {
	cur := CurrentLLMSettings()
	llmMu.RLock()
	defer llmMu.RUnlock()
	return LLMStatus{
		Configured: cur.APIKey != "",
		Provider:   cur.Provider,
		Model:      cur.Model,
		BaseURL:    cur.BaseURL,
		LastSource: llmLastSource,
		LastError:  llmLastError,
	}
}

type LLMClient struct {
	provider string
	apiKey   string
	baseURL  string
	model    string
	client   *http.Client
}

func NewLLMClient() *LLMClient {
	s := CurrentLLMSettings()
	return &LLMClient{
		provider: s.Provider,
		apiKey:   s.APIKey,
		baseURL:  s.BaseURL,
		model:    s.Model,
		client:   &http.Client{Timeout: 25 * time.Second},
	}
}

func (c *LLMClient) IsConfigured() bool {
	return c.apiKey != ""
}

type LLMActivityResponse struct {
	Title                 string `json:"title"`
	Content               string `json:"content"`
	PersonalizationReason string `json:"personalization_reason"`
}

const llmSystemPrompt = `You are an expert instructional designer for an adaptive learning platform.
Rewrite the given learning activity to better engage a student based on their current engagement state and the intervention strategy.
Return ONLY a JSON object with exactly these keys:
- "title": string, max 140 chars
- "content": string, the adapted activity content, max 4000 chars
- "personalization_reason": string, one sentence explaining why this adaptation fits the student`

func (c *LLMClient) complete(ctx context.Context, system, user string) (string, error) {
	var url string
	var body map[string]any
	headers := map[string]string{"Content-Type": "application/json"}

	switch c.provider {
	case "anthropic":
		url = c.baseURL + "/v1/messages"
		headers["x-api-key"] = c.apiKey
		headers["anthropic-version"] = "2023-06-01"
		body = map[string]any{
			"model":      c.model,
			"max_tokens": 1024,
			"system":     system,
			"messages":   []map[string]string{{"role": "user", "content": user}},
		}
	case "gemini":
		url = fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
		body = map[string]any{
			"systemInstruction": map[string]any{
				"parts": []map[string]string{{"text": system}},
			},
			"contents": []map[string]any{
				{"role": "user", "parts": []map[string]string{{"text": user}}},
			},
			"generationConfig": map[string]any{
				"temperature":     0.3,
				"maxOutputTokens": 1024,
			},
		}
	default: // openai-compatible
		url = c.baseURL + "/chat/completions"
		headers["Authorization"] = "Bearer " + c.apiKey
		body = map[string]any{
			"model": c.model,
			"messages": []map[string]string{
				{"role": "system", "content": system},
				{"role": "user", "content": user},
			},
			"temperature": 0.3,
			"max_tokens":  1024,
		}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s", strings.ReplaceAll(err.Error(), c.apiKey, "***"))
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK {
		msg := strings.ReplaceAll(string(data), c.apiKey, "***")
		if len(msg) > 200 {
			msg = msg[:200]
		}
		return "", fmt.Errorf("llm request failed (status %d): %s", resp.StatusCode, msg)
	}

	if c.provider == "anthropic" {
		var r struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(data, &r); err != nil {
			return "", err
		}
		var sb strings.Builder
		for _, b := range r.Content {
			if b.Type == "text" {
				sb.WriteString(b.Text)
			}
		}
		if sb.Len() == 0 {
			return "", fmt.Errorf("empty response from LLM")
		}
		return sb.String(), nil
	}

	if c.provider == "gemini" {
		var r struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}
		if err := json.Unmarshal(data, &r); err != nil {
			return "", err
		}
		if len(r.Candidates) == 0 {
			return "", fmt.Errorf("empty response from Gemini")
		}
		var sb strings.Builder
		for _, p := range r.Candidates[0].Content.Parts {
			sb.WriteString(p.Text)
		}
		if sb.Len() == 0 {
			return "", fmt.Errorf("empty text from Gemini")
		}
		return sb.String(), nil
	}

	var r struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return "", err
	}
	if len(r.Choices) == 0 {
		return "", fmt.Errorf("empty response choices from LLM")
	}
	return r.Choices[0].Message.Content, nil
}

func validateLLMActivity(r *LLMActivityResponse) error {
	r.Title = strings.TrimSpace(r.Title)
	r.Content = strings.TrimSpace(r.Content)
	r.PersonalizationReason = strings.TrimSpace(r.PersonalizationReason)
	switch {
	case r.Title == "" || len(r.Title) > 140:
		return fmt.Errorf("invalid title from LLM")
	case r.Content == "" || len(r.Content) > 4000:
		return fmt.Errorf("invalid content from LLM")
	}
	return nil
}

func (c *LLMClient) GenerateAdaptedActivity(
	ctx context.Context,
	orig models.Activity,
	interventionType string,
	interventionReasons []string,
	studentState string,
) (*LLMActivityResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("LLM API key not configured")
	}

	user := fmt.Sprintf(
		"Original Activity:\nTitle: %s\nTopic: %s\nDifficulty: %s\nContent:\n%s\nExpected Answer: %q\n\nStudent engagement state: %s\nIntervention: %s (reason: %s)\n\nReturn the adapted activity as JSON.",
		orig.Title, orig.Topic, orig.Difficulty, orig.Content, orig.CorrectAnswer,
		studentState, interventionType, strings.Join(interventionReasons, ", "),
	)

	raw, err := c.complete(ctx, llmSystemPrompt, user)
	if err != nil {
		return nil, err
	}
	if i, j := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); i >= 0 && j > i {
		raw = raw[i : j+1]
	}

	var res LLMActivityResponse
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, fmt.Errorf("failed to parse structured JSON: %w", err)
	}
	if err := validateLLMActivity(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
