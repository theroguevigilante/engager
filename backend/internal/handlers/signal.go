package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adaptive-learning/backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) RecordSignals(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "id")
	if attemptID == "" {
		http.Error(w, "missing attempt id parameter", http.StatusBadRequest)
		return
	}

	var batch models.TelemetryBatch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.IngestSignals(r.Context(), attemptID, batch); err != nil {
		http.Error(w, "failed to record signals: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "signals_recorded",
	})
}
