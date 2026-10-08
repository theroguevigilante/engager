package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/adaptive-learning/backend/internal/models"
)

func (r *Repository) SaveAdaptationLoop(ctx context.Context, loop models.AdaptationLoop) error {
	reasonsBytes, err := json.Marshal(loop.EffectivenessReasons)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO adaptation_loops (
			id, student_id, original_attempt_id, adapted_attempt_id, assessment_id, intervention_id,
			personalized_activity_id, pre_state, post_state, pre_engagement_score, post_engagement_score,
			engagement_delta, effectiveness, effectiveness_reasons_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		loop.ID,
		loop.StudentID,
		loop.OriginalAttemptID,
		loop.AdaptedAttemptID,
		loop.AssessmentID,
		loop.InterventionID,
		loop.PersonalizedActivityID,
		string(loop.PreState),
		string(loop.PostState),
		loop.PreEngagementScore,
		loop.PostEngagementScore,
		loop.EngagementDelta,
		string(loop.Effectiveness),
		string(reasonsBytes),
		loop.CreatedAt,
	)
	return err
}

func (r *Repository) GetLatestAdaptationLoopByOriginalAttempt(ctx context.Context, originalAttemptID string) (*models.AdaptationLoop, error) {
	query := `
		SELECT id, student_id, original_attempt_id, adapted_attempt_id, assessment_id, intervention_id,
		       personalized_activity_id, pre_state, post_state, pre_engagement_score, post_engagement_score,
		       engagement_delta, effectiveness, effectiveness_reasons_json, created_at
		FROM adaptation_loops
		WHERE original_attempt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var (
		loop        models.AdaptationLoop
		preState    string
		postState   string
		effStr      string
		reasonsJSON string
	)

	err := r.db.QueryRowContext(ctx, query, originalAttemptID).Scan(
		&loop.ID,
		&loop.StudentID,
		&loop.OriginalAttemptID,
		&loop.AdaptedAttemptID,
		&loop.AssessmentID,
		&loop.InterventionID,
		&loop.PersonalizedActivityID,
		&preState,
		&postState,
		&loop.PreEngagementScore,
		&loop.PostEngagementScore,
		&loop.EngagementDelta,
		&effStr,
		&reasonsJSON,
		&loop.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	loop.PreState = models.EngagementState(preState)
	loop.PostState = models.EngagementState(postState)
	loop.Effectiveness = models.LoopEffectiveness(effStr)
	if err := json.Unmarshal([]byte(reasonsJSON), &loop.EffectivenessReasons); err != nil {
		return nil, err
	}
	return &loop, nil
}

func (r *Repository) GetAdaptedAttemptByParentID(ctx context.Context, parentAttemptID string) (*models.ActivityAttempt, error) {
	query := `
		SELECT id, student_id, activity_id, parent_attempt_id, started_at, completed_at, duration, status, attempt_count, created_at
		FROM activity_attempts
		WHERE parent_attempt_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var (
		att         models.ActivityAttempt
		completedAt sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, query, parentAttemptID).Scan(
		&att.ID,
		&att.StudentID,
		&att.ActivityID,
		&att.ParentAttemptID,
		&att.StartedAt,
		&completedAt,
		&att.Duration,
		&att.Status,
		&att.AttemptCount,
		&att.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if completedAt.Valid {
		att.CompletedAt = &completedAt.Time
	}
	return &att, nil
}

const loopColumns = `id, student_id, original_attempt_id, adapted_attempt_id, assessment_id, intervention_id,
	personalized_activity_id, pre_state, post_state, pre_engagement_score, post_engagement_score,
	engagement_delta, effectiveness, effectiveness_reasons_json, created_at`

func scanLoop(sc interface{ Scan(...any) error }) (models.AdaptationLoop, error) {
	var (
		loop                                  models.AdaptationLoop
		preState, postState, effStr, reasonsJ string
	)
	if err := sc.Scan(&loop.ID, &loop.StudentID, &loop.OriginalAttemptID, &loop.AdaptedAttemptID,
		&loop.AssessmentID, &loop.InterventionID, &loop.PersonalizedActivityID, &preState, &postState,
		&loop.PreEngagementScore, &loop.PostEngagementScore, &loop.EngagementDelta, &effStr, &reasonsJ,
		&loop.CreatedAt); err != nil {
		return loop, err
	}
	loop.PreState = models.EngagementState(preState)
	loop.PostState = models.EngagementState(postState)
	loop.Effectiveness = models.LoopEffectiveness(effStr)
	loop.EffectivenessReasons = []string{}
	if err := json.Unmarshal([]byte(reasonsJ), &loop.EffectivenessReasons); err != nil {
		return loop, err
	}
	return loop, nil
}

// ListAdaptationLoops returns closed loops for a student, newest first.
func (r *Repository) ListAdaptationLoops(ctx context.Context, studentID string, limit int) ([]models.AdaptationLoop, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT "+loopColumns+" FROM adaptation_loops WHERE student_id = ? ORDER BY created_at DESC LIMIT ?",
		studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loops := []models.AdaptationLoop{}
	for rows.Next() {
		l, err := scanLoop(rows)
		if err != nil {
			return nil, err
		}
		loops = append(loops, l)
	}
	return loops, rows.Err()
}
