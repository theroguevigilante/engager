package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) EvaluateAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "id")
	if attemptID == "" {
		http.Error(w, "missing attempt id parameter", http.StatusBadRequest)
		return
	}

	result, err := h.svc.EvaluateAttempt(r.Context(), attemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAssessment(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "id")
	if attemptID == "" {
		http.Error(w, "missing attempt id parameter", http.StatusBadRequest)
		return
	}

	result, err := h.svc.GetAssessment(r.Context(), attemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "assessment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
