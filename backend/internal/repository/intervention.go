package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adaptive-learning/backend/internal/models"
)

func (r *Repository) SaveIntervention(ctx context.Context, inter models.AttemptIntervention) error {
	query := `
		INSERT INTO attempt_interventions (id, attempt_id, assessment_id, intervention_type, reason, confidence, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		inter.ID,
		inter.AttemptID,
		inter.AssessmentID,
		string(inter.InterventionType),
		inter.Reason,
		inter.Confidence,
		inter.CreatedAt,
	)
	return err
}

func (r *Repository) GetLatestInterventionByAttemptID(ctx context.Context, attemptID string) (*models.AttemptIntervention, error) {
	query := `
		SELECT id, attempt_id, assessment_id, intervention_type, reason, confidence, created_at
		FROM attempt_interventions
		WHERE attempt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var (
		inter   models.AttemptIntervention
		typeStr string
	)

	err := r.db.QueryRowContext(ctx, query, attemptID).Scan(
		&inter.ID,
		&inter.AttemptID,
		&inter.AssessmentID,
		&typeStr,
		&inter.Reason,
		&inter.Confidence,
		&inter.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	inter.InterventionType = models.InterventionType(typeStr)
	return &inter, nil
}
