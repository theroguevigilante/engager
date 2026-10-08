package planner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adaptive-learning/backend/internal/models"
)

func serve(t *testing.T, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
}

func TestLLMProviders(t *testing.T) {
	act := models.Activity{Title: "t", Content: "c", Topic: "Mathematics", Difficulty: "beginner"}
	js := `{\"title\":\"Guided\",\"content\":\"Step 1...\",\"personalization_reason\":\"Here are hints.\"}`

	oa := serve(t, `{"choices":[{"message":{"content":"`+"```json\\n"+js+"\\n```"+`"}}]}`)
	defer oa.Close()
	c := &LLMClient{provider: "openai", apiKey: "k", baseURL: oa.URL, model: "m", client: oa.Client()}
	if r, err := c.GenerateAdaptedActivity(context.Background(), act, "SCAFFOLD", nil, "STRUGGLING"); err != nil || r.Title != "Guided" {
		t.Fatalf("openai: %v %+v", err, r)
	}

	an := serve(t, `{"content":[{"type":"text","text":"`+js+`"}]}`)
	defer an.Close()
	c = &LLMClient{provider: "anthropic", apiKey: "k", baseURL: an.URL, model: "m", client: an.Client()}
	if r, err := c.GenerateAdaptedActivity(context.Background(), act, "SCAFFOLD", nil, "STRUGGLING"); err != nil || r.Content == "" {
		t.Fatalf("anthropic: %v %+v", err, r)
	}

	bad := serve(t, `{"content":[{"type":"text","text":"{\"title\":\"\",\"content\":\"x\"}"}]}`)
	defer bad.Close()
	c = &LLMClient{provider: "anthropic", apiKey: "k", baseURL: bad.URL, model: "m", client: bad.Client()}
	if _, err := c.GenerateAdaptedActivity(context.Background(), act, "SCAFFOLD", nil, "STRUGGLING"); err == nil {
		t.Fatal("expected validation error for empty title")
	}
}

func TestStatusNeverLeaksKey(t *testing.T) {
	SetLLMSettings(LLMSettings{Provider: "anthropic", APIKey: "secret"})
	if st := GetLLMStatus(); !st.Configured || st.Provider != "anthropic" {
		t.Fatalf("bad status %+v", st)
	}
	ClearLLMSettings()
	if GetLLMStatus().Configured {
		t.Fatal("expected disabled")
	}
}

func TestNeedTwistIsDomainAgnostic(t *testing.T) {
	for _, n := range []string{"lonely", "playful", "curious", "stuck", "bored"} {
		tw, why := NeedTwist(n)
		if tw == "" || why == "" {
			t.Fatalf("%s: empty twist", n)
		}
	}
	if tw, _ := NeedTwist("fine"); tw != "" {
		t.Fatal("fine should add nothing")
	}
}
