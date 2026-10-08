package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
	"github.com/adaptive-learning/backend/internal/planner"
)

func (s *Service) PersonalizeActivityForAttempt(ctx context.Context, attemptID string) (*models.PersonalizedActivity, error) {
	attempt, err := s.repo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("attempt not found: %w", err)
	}

	assessment, err := s.repo.GetLatestAssessmentByAttemptID(ctx, attemptID)
	if err != nil || assessment == nil {
		return nil, fmt.Errorf("assessment not found for attempt %s", attemptID)
	}

	interv, err := s.repo.GetLatestInterventionByAttemptID(ctx, attemptID)
	if err != nil || interv == nil {
		return nil, fmt.Errorf("intervention not found for attempt %s", attemptID)
	}

	currentActivity, err := s.repo.GetActivityByID(ctx, attempt.ActivityID)
	if err != nil {
		return nil, fmt.Errorf("current activity not found: %w", err)
	}

	// Prepare reasons slice from AttemptIntervention.Reason
	var reasons []string
	if interv.Reason != "" {
		reasons = []string{interv.Reason}
	}

	// Build the InterventionSpec expected by the updated planner
	spec := planner.InterventionSpec{
		InterventionType: interv.InterventionType,
		Reasons:          reasons,
		StudentState:     string(assessment.State),
		Need:             s.LatestNeed(attempt.StudentID),
	}

	// Instantiate the ActivityPlanner (with LLM client & domain fallback)
	actPlanner := planner.NewActivityPlanner(s.repo)
	planned, err := actPlanner.PlanPersonalizedActivity(ctx, attempt, spec, *currentActivity)
	if err != nil {
		return nil, fmt.Errorf("failed planning personalized activity: %w", err)
	}

	planned.AttemptID = attemptID
	planned.InterventionID = interv.ID
	planned.CreatedAt = time.Now().UTC()

	if err := s.repo.SavePersonalizedActivity(ctx, *planned); err != nil {
		return nil, fmt.Errorf("failed saving personalized activity: %w", err)
	}

	return planned, nil
}

func (s *Service) GetPersonalizedActivity(ctx context.Context, attemptID string) (*models.PersonalizedActivity, error) {
	return s.repo.GetLatestPersonalizedActivityByAttemptID(ctx, attemptID)
}
