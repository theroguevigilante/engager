package main

import (
	"log"

	"github.com/adaptive-learning/backend/internal/repository"
	"github.com/adaptive-learning/backend/internal/service"
	"github.com/adaptive-learning/backend/migrations"
)

// Adds the starter activity library (even if activities already exist).
func main() {
	repo, err := repository.NewRepository("./engager.db")
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	if err := repo.Migrate(migrations.FS, "."); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	if err := service.NewService(repo).SeedDefaults(true); err != nil {
		log.Fatalf("Seed error: %v", err)
	}
	log.Println("Database seeded successfully.")
}
