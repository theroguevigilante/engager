package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adaptive-learning/backend/internal/handlers"
	"github.com/adaptive-learning/backend/internal/repository"
	"github.com/adaptive-learning/backend/internal/service"
	"github.com/adaptive-learning/backend/migrations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialize local SQLite database
	repo, err := repository.NewRepository("./engager.db")
	if err != nil {
		log.Fatalf("Failed to initialize sqlite repository: %v", err)
	}

	// Apply all database migrations (idempotent)
	if err := repo.Migrate(migrations.FS, "."); err != nil {
		log.Fatalf("Failed to apply schema migration: %v", err)
	}

	// Wire service and handler layers
	svc := service.NewService(repo)
	if err := svc.SeedDefaults(false); err != nil {
		log.Printf("seeding skipped: %v", err)
	}
	h := handlers.NewHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Development CORS configuration for Vite frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "ok",
			})
		})

		// Activities & session tracking
		r.Get("/activities", h.GetActivities)
		r.Post("/activities", h.CreateActivity)
		r.Post("/activities/{id}/start", h.StartActivity)
		r.Post("/activities/{id}/complete", h.CompleteActivity)
		r.Post("/activities/{id}/abandon", h.AbandonActivity)

		// Student telemetry & history
		r.Get("/students/{id}/activity-history", h.GetStudentHistory)
		r.Post("/attempts/{id}/signals", h.RecordSignals)
		r.Post("/attempts/{id}/evaluate", h.EvaluateAttempt)
		r.Get("/attempts/{id}/assessment", h.GetAssessment)
		r.Post("/attempts/{id}/intervene", h.SelectIntervention)
		r.Get("/attempts/{id}/intervention", h.GetIntervention)
		r.Post("/attempts/{id}/personalize", h.PersonalizeAttempt)
		r.Get("/attempts/{id}/personalized-activity", h.GetPersonalizedActivity)
		r.Post("/attempts/{id}/close-loop", h.CloseLoop)
		r.Get("/attempts/{id}/adaptation-loop", h.GetLoop)

		// Adaptation history
		r.Get("/loops", h.ListLoops)
		r.Get("/loops/latest", h.GetLatestLoop)
		r.Post("/demo/run", h.RunDemo)

		// Student check-in (self-reported need) and anonymous community pulse
		r.Post("/students/{id}/checkin", h.PostCheckin)
		r.Get("/community/pulse", h.GetPulse)

		// Optional LLM configuration (content adaptation only)
		r.Get("/settings/llm", h.GetLLMSettings)
		r.Put("/settings/llm", h.PutLLMSettings)
		r.Delete("/settings/llm", h.DeleteLLMSettings)

		r.Put("/activities/{id}", h.UpdateActivity)
	})

	log.Println("Backend server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
