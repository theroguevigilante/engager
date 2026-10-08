package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/adaptive-learning/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.svc.GetActivities()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if activities == nil {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, activities)
}

func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var body models.Activity
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	act, err := h.svc.CreateActivity(body)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidActivity) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	act.CorrectAnswer = ""
	writeJSON(w, http.StatusCreated, act)
}

func (h *Handler) StartActivity(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "id")
	var body struct {
		StudentID       string  `json:"student_id"`
		ParentAttemptID *string `json:"parent_attempt_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.StudentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "student_id is required"})
		return
	}

	attempt, err := h.svc.StartActivity(activityID, body.StudentID, body.ParentAttemptID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, attempt)
}

func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var act models.Activity
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}
	act.ID = id

	if err := h.svc.UpdateActivity(act); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrInvalidActivity):
			status = http.StatusBadRequest
		case errors.Is(err, sql.ErrNoRows):
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	act.CorrectAnswer = ""
	writeJSON(w, http.StatusOK, act)
}

func (h *Handler) CompleteActivity(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "id")
	var body struct {
		AttemptID string `json:"attempt_id"`
		StudentID string `json:"student_id"`
		Answer    string `json:"answer"`
		Attempts  int    `json:"attempts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	attempt, correct, err := h.svc.CompleteActivity(activityID, body.AttemptID, body.StudentID, body.Answer, body.Attempts)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": correct,
		"attempt": attempt,
	})
}

func (h *Handler) AbandonActivity(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "id")
	var body struct {
		AttemptID string `json:"attempt_id"`
		StudentID string `json:"student_id"`
		Attempts  int    `json:"attempts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	attempt, err := h.svc.AbandonActivity(activityID, body.AttemptID, body.StudentID, body.Attempts)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, attempt)
}

func (h *Handler) GetStudentHistory(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "id")
	history, err := h.svc.GetStudentHistory(studentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if history == nil {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, history)
}
