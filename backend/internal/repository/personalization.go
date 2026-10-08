package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adaptive-learning/backend/internal/models"
)

func (r *Repository) SavePersonalizedActivity(ctx context.Context, p models.PersonalizedActivity) error {
	query := `
		INSERT INTO personalized_activities (
			id, attempt_id, source_activity_id, intervention_id, title, content, difficulty, personalization_reason, is_generated, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.AttemptID,
		p.SourceActivityID,
		p.InterventionID,
		p.Title,
		p.Content,
		p.Difficulty,
		p.PersonalizationReason,
		p.IsGenerated,
		p.CreatedAt,
	)
	return err
}

func (r *Repository) GetLatestPersonalizedActivityByAttemptID(ctx context.Context, attemptID string) (*models.PersonalizedActivity, error) {
	query := `
		SELECT id, attempt_id, source_activity_id, intervention_id, title, content, difficulty, personalization_reason, is_generated, created_at
		FROM personalized_activities
		WHERE attempt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var p models.PersonalizedActivity
	err := r.db.QueryRowContext(ctx, query, attemptID).Scan(
		&p.ID,
		&p.AttemptID,
		&p.SourceActivityID,
		&p.InterventionID,
		&p.Title,
		&p.Content,
		&p.Difficulty,
		&p.PersonalizationReason,
		&p.IsGenerated,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetAllActivities(ctx context.Context) ([]models.Activity, error) {
	// Query handles the base fields; if your table does not have topic/difficulty columns yet:
	query := `SELECT id, title, description, content, correct_answer, created_at FROM activities`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.Content, &a.CorrectAnswer, &a.CreatedAt); err != nil {
			return nil, err
		}
		// Fallbacks in case columns aren't in SQLite table
		if a.Difficulty == "" {
			a.Difficulty = "beginner"
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *Repository) GetActivityByID(ctx context.Context, id string) (*models.Activity, error) {
	query := `SELECT id, title, description, content, correct_answer, created_at FROM activities WHERE id = ?`
	var a models.Activity
	err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Title, &a.Description, &a.Content, &a.CorrectAnswer, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	if a.Difficulty == "" {
		a.Difficulty = "beginner"
	}
	return &a, nil
}
