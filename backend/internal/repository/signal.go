package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
	"github.com/google/uuid"
)

func (r *Repository) SaveSignalBatch(ctx context.Context, attemptID string, batch models.TelemetryBatch) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO engagement_signals (id, student_id, attempt_id, signal_type, signal_value, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC()

	// Map batch fields to key-value signal rows
	signals := []struct {
		sType string
		sVal  string
	}{
		{"session_duration", strconv.Itoa(batch.SessionDuration)},
		{"active_duration", strconv.Itoa(batch.ActiveDuration)},
		{"idle_duration", strconv.Itoa(batch.IdleDuration)},
		{"tab_switches", strconv.Itoa(batch.TabSwitches)},
		{"interaction_count", strconv.Itoa(batch.InteractionCount)},
	}

	for _, s := range signals {
		id := uuid.New().String()
		_, err := stmt.ExecContext(ctx, id, batch.StudentID, attemptID, s.sType, s.sVal, now)
		if err != nil {
			return fmt.Errorf("failed inserting signal %s: %w", s.sType, err)
		}
	}

	return tx.Commit()
}
