package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adaptive-learning/backend/internal/assessment"
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/google/uuid"
)

func (s *Service) EvaluateAttempt(ctx context.Context, attemptID string) (*models.EngagementAssessment, error) {
	attempt, err := s.repo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("attempt not found: %w", err)
	}

	signals, err := s.repo.GetSignalsByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("failed fetching signals: %w", err)
	}

	metrics := assessment.ComputeMetrics(signals, *attempt)
	state, confidence, reasons := assessment.Classify(metrics)

	eval := models.EngagementAssessment{
		ID:         uuid.New().String(),
		AttemptID:  attemptID,
		State:      state,
		Confidence: confidence,
		Metrics:    metrics,
		Reasons:    reasons,
		AssessedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveAssessment(ctx, eval); err != nil {
		return nil, fmt.Errorf("failed persisting assessment: %w", err)
	}

	return &eval, nil
}

func (s *Service) GetAssessment(ctx context.Context, attemptID string) (*models.EngagementAssessment, error) {
	return s.repo.GetLatestAssessmentByAttemptID(ctx, attemptID)
}
