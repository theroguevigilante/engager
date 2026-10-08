package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CloseLoop(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "id")
	if attemptID == "" {
		http.Error(w, "missing attempt id parameter", http.StatusBadRequest)
		return
	}

	result, err := h.svc.CloseAdaptationLoop(r.Context(), attemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetLoop(w http.ResponseWriter, r *http.Request) {
	attemptID := chi.URLParam(r, "id")
	if attemptID == "" {
		http.Error(w, "missing attempt id parameter", http.StatusBadRequest)
		return
	}

	result, err := h.svc.GetAdaptationLoop(r.Context(), attemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "adaptation loop record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

const defaultStudentID = "student-demo-101"

// ListLoops returns closed adaptation loops for the student (?student_id=, demo student by default).
func (h *Handler) ListLoops(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		studentID = defaultStudentID
	}
	loops, err := h.svc.ListAdaptationLoops(r.Context(), studentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, loops)
}

// GetLatestLoop returns the most recently closed loop for the student.
func (h *Handler) GetLatestLoop(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		studentID = defaultStudentID
	}
	loops, err := h.svc.ListAdaptationLoops(r.Context(), studentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if len(loops) == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no adaptation loops yet"})
		return
	}
	writeJSON(w, http.StatusOK, loops[0])
}

// RunDemo executes a scripted end-to-end adaptation loop for a chosen activity.
func (h *Handler) RunDemo(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ActivityID string `json:"activity_id"`
		StudentID  string `json:"student_id"`
		Scenario   string `json:"scenario"`
		Need       string `json:"need"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ActivityID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "activity_id is required"})
		return
	}
	if body.StudentID == "" {
		body.StudentID = defaultStudentID
	}
	res, err := h.svc.RunDemoLoop(r.Context(), body.ActivityID, body.StudentID, body.Scenario, body.Need)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
