package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adaptive-learning/backend/internal/intervention"
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/google/uuid"
)

func (s *Service) SelectAndSaveIntervention(ctx context.Context, attemptID string) (*models.AttemptIntervention, error) {
	ass, err := s.repo.GetLatestAssessmentByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("failed fetching assessment: %w", err)
	}
	if ass == nil {
		return nil, fmt.Errorf("no engagement assessment found for attempt %s", attemptID)
	}

	att, err := s.repo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("failed fetching attempt: %w", err)
	}

	iType, reason, conf := intervention.SelectIntervention(ass, att)
	metadata := intervention.Catalog[iType]

	item := models.AttemptIntervention{
		ID:               uuid.New().String(),
		AttemptID:        attemptID,
		AssessmentID:     ass.ID,
		InterventionType: iType,
		Reason:           reason,
		Confidence:       conf,
		Metadata:         metadata,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.repo.SaveIntervention(ctx, item); err != nil {
		return nil, fmt.Errorf("failed saving intervention: %w", err)
	}

	return &item, nil
}

func (s *Service) GetIntervention(ctx context.Context, attemptID string) (*models.AttemptIntervention, error) {
	inter, err := s.repo.GetLatestInterventionByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if inter == nil {
		return nil, nil
	}
	inter.Metadata = intervention.Catalog[inter.InterventionType]
	return inter, nil
}
