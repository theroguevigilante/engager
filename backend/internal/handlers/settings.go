package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adaptive-learning/backend/internal/planner"
	"github.com/adaptive-learning/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

// GetLLMSettings reports LLM configuration status. The API key is never returned.
func (h *Handler) GetLLMSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, planner.GetLLMStatus())
}

// PutLLMSettings stores provider/key/model in memory (not persisted to disk).
func (h *Handler) PutLLMSettings(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Provider string `json:"provider"`
		APIKey   string `json:"api_key"`
		Model    string `json:"model"`
		BaseURL  string `json:"base_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}
	planner.SetLLMSettings(planner.LLMSettings{
		Provider: body.Provider, APIKey: body.APIKey, Model: body.Model, BaseURL: body.BaseURL,
	})
	writeJSON(w, http.StatusOK, planner.GetLLMStatus())
}

// DeleteLLMSettings disables LLM adaptation (deterministic rules only).
func (h *Handler) DeleteLLMSettings(w http.ResponseWriter, r *http.Request) {
	planner.ClearLLMSettings()
	writeJSON(w, http.StatusOK, planner.GetLLMStatus())
}

// PostCheckin records how the student says they feel ("stuck", "lonely", ...).
func (h *Handler) PostCheckin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Need string `json:"need"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}
	if err := h.svc.RecordCheckin(chi.URLParam(r, "id"), body.Need); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidNeed) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"need": body.Need})
}

// GetPulse returns anonymous counts of self-reported needs in the last 24 hours.
func (h *Handler) GetPulse(w http.ResponseWriter, r *http.Request) {
	counts, err := h.svc.CommunityPulse()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"counts": counts})
}
