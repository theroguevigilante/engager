package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adaptive-learning/backend/internal/loop"
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/google/uuid"
)

func (s *Service) CloseAdaptationLoop(ctx context.Context, originalAttemptID string) (*models.AdaptationLoop, error) {
	// 1. Fetch original attempt
	origAttempt, err := s.repo.GetAttempt(originalAttemptID)
	if err != nil {
		return nil, fmt.Errorf("original attempt not found: %w", err)
	}

	// 2. Fetch original attempt's assessment
	preAssessment, err := s.repo.GetLatestAssessmentByAttemptID(ctx, originalAttemptID)
	if err != nil || preAssessment == nil {
		return nil, fmt.Errorf("pre-assessment not found for attempt %s", originalAttemptID)
	}

	// 3. Fetch intervention selected for original attempt
	intervention, err := s.repo.GetLatestInterventionByAttemptID(ctx, originalAttemptID)
	if err != nil || intervention == nil {
		return nil, fmt.Errorf("intervention not found for attempt %s", originalAttemptID)
	}

	// 4. Fetch the personalized activity generated
	personalizedAct, err := s.repo.GetLatestPersonalizedActivityByAttemptID(ctx, originalAttemptID)
	if err != nil || personalizedAct == nil {
		return nil, fmt.Errorf("personalized activity not found for attempt %s", originalAttemptID)
	}

	// 5. Fetch adapted attempt linked via parent_attempt_id
	adaptedAttempt, err := s.repo.GetAdaptedAttemptByParentID(ctx, originalAttemptID)
	if err != nil || adaptedAttempt == nil {
		return nil, fmt.Errorf("adapted attempt not found for parent attempt %s: %w", originalAttemptID, err)
	}

	// 6. Fetch post-adaptation assessment (or evaluate if not evaluated yet)
	postAssessment, err := s.repo.GetLatestAssessmentByAttemptID(ctx, adaptedAttempt.ID)
	if err != nil || postAssessment == nil {
		postAssessment, err = s.EvaluateAttempt(ctx, adaptedAttempt.ID)
		if err != nil {
			return nil, fmt.Errorf("failed evaluating adapted attempt %s: %w", adaptedAttempt.ID, err)
		}
	}

	// 7. Calculate before vs after engagement & determine effectiveness
	effectiveness, preScore, postScore, delta, reasons := loop.EvaluateLoop(
		preAssessment,
		postAssessment,
		adaptedAttempt,
		intervention.InterventionType,
	)

	record := models.AdaptationLoop{
		ID:                     uuid.New().String(),
		StudentID:              origAttempt.StudentID,
		OriginalAttemptID:      originalAttemptID,
		AdaptedAttemptID:       adaptedAttempt.ID,
		AssessmentID:           preAssessment.ID,
		InterventionID:         intervention.ID,
		PersonalizedActivityID: personalizedAct.ID,
		PreState:               preAssessment.State,
		PostState:              postAssessment.State,
		PreEngagementScore:     preScore,
		PostEngagementScore:    postScore,
		EngagementDelta:        delta,
		Effectiveness:          effectiveness,
		EffectivenessReasons:   reasons,
		CreatedAt:              time.Now().UTC(),
	}

	// 8. Persist the loop result
	if err := s.repo.SaveAdaptationLoop(ctx, record); err != nil {
		return nil, fmt.Errorf("failed saving adaptation loop: %w", err)
	}

	return &record, nil
}

func (s *Service) GetAdaptationLoop(ctx context.Context, originalAttemptID string) (*models.AdaptationLoop, error) {
	return s.repo.GetLatestAdaptationLoopByOriginalAttempt(ctx, originalAttemptID)
}

func (s *Service) ListAdaptationLoops(ctx context.Context, studentID string) ([]models.AdaptationLoop, error) {
	return s.repo.ListAdaptationLoops(ctx, studentID, 50)
}
