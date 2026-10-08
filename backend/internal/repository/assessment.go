package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/adaptive-learning/backend/internal/models"
)

func (r *Repository) SaveAssessment(ctx context.Context, a models.EngagementAssessment) error {
	metricsBytes, err := json.Marshal(a.Metrics)
	if err != nil {
		return err
	}
	reasonsBytes, err := json.Marshal(a.Reasons)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO engagement_assessments (id, attempt_id, state, confidence, metrics_json, reasons_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		a.ID,
		a.AttemptID,
		string(a.State),
		a.Confidence,
		string(metricsBytes),
		string(reasonsBytes),
		a.AssessedAt,
	)
	return err
}

func (r *Repository) GetLatestAssessmentByAttemptID(ctx context.Context, attemptID string) (*models.EngagementAssessment, error) {
	query := `
		SELECT id, attempt_id, state, confidence, metrics_json, reasons_json, created_at
		FROM engagement_assessments
		WHERE attempt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var (
		a           models.EngagementAssessment
		stateStr    string
		metricsJSON string
		reasonsJSON string
	)

	err := r.db.QueryRowContext(ctx, query, attemptID).Scan(
		&a.ID,
		&a.AttemptID,
		&stateStr,
		&a.Confidence,
		&metricsJSON,
		&reasonsJSON,
		&a.AssessedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	a.State = models.EngagementState(stateStr)
	if err := json.Unmarshal([]byte(metricsJSON), &a.Metrics); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(reasonsJSON), &a.Reasons); err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *Repository) GetAttemptByID(ctx context.Context, id string) (*models.ActivityAttempt, error) {
	query := `
		SELECT id, student_id, activity_id, started_at, completed_at, duration, status, attempt_count, created_at
		FROM activity_attempts
		WHERE id = ?
	`
	var (
		att         models.ActivityAttempt
		completedAt sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&att.ID,
		&att.StudentID,
		&att.ActivityID,
		&att.StartedAt,
		&completedAt,
		&att.Duration,
		&att.Status,
		&att.AttemptCount,
		&att.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		att.CompletedAt = &completedAt.Time
	}
	return &att, nil
}

func (r *Repository) GetSignalsByAttemptID(ctx context.Context, attemptID string) ([]models.EngagementSignal, error) {
	query := `
		SELECT id, student_id, attempt_id, signal_type, signal_value, created_at
		FROM engagement_signals
		WHERE attempt_id = ?
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []models.EngagementSignal
	for rows.Next() {
		var s models.EngagementSignal
		if err := rows.Scan(&s.ID, &s.StudentID, &s.AttemptID, &s.SignalType, &s.SignalValue, &s.CreatedAt); err != nil {
			return nil, err
		}
		signals = append(signals, s)
	}
	return signals, rows.Err()
}
