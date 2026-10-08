package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
	"github.com/adaptive-learning/backend/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateActivity(act models.Activity) error {
	if err := normalizeActivity(&act); err != nil {
		return err
	}
	return s.repo.UpdateActivity(act)
}

func (s *Service) GetActivities() ([]models.Activity, error) {
	return s.repo.ListActivities()
}

var ErrInvalidActivity = errors.New("title and content are required")

func normalizeActivity(a *models.Activity) error {
	a.Title = strings.TrimSpace(a.Title)
	a.Content = strings.TrimSpace(a.Content)
	if a.Title == "" || a.Content == "" {
		return ErrInvalidActivity
	}
	if strings.TrimSpace(a.Topic) == "" {
		a.Topic = "General"
	}
	switch a.Difficulty {
	case "beginner", "intermediate", "advanced":
	default:
		a.Difficulty = "intermediate"
	}
	if a.EstimatedDurationMinutes <= 0 {
		a.EstimatedDurationMinutes = 5
	}
	return nil
}

func (s *Service) CreateActivity(act models.Activity) (*models.Activity, error) {
	if err := normalizeActivity(&act); err != nil {
		return nil, err
	}
	act.ID = uuid.New().String()
	if err := s.repo.CreateActivity(act); err != nil {
		return nil, err
	}
	return &act, nil
}

func (s *Service) StartActivity(activityID, studentID string, parentAttemptID *string) (*models.ActivityAttempt, error) {
	_ = s.repo.EnsureStudent(studentID, "Demo Student")

	att := models.ActivityAttempt{
		ID:              uuid.New().String(),
		StudentID:       studentID,
		ActivityID:      activityID,
		ParentAttemptID: parentAttemptID,
		StartedAt:       time.Now(),
		Status:          "in_progress",
		AttemptCount:    1,
	}

	if err := s.repo.CreateAttempt(att); err != nil {
		return nil, err
	}
	return &att, nil
}

func (s *Service) CompleteActivity(activityID string, attemptID string, studentID string, answer string, attempts int) (*models.ActivityAttempt, bool, error) {
	att, err := s.repo.GetAttempt(attemptID)
	if err != nil {
		return nil, false, err
	}

	act, err := s.repo.GetActivity(activityID)
	if err != nil {
		return nil, false, err
	}

	expected := strings.TrimSpace(act.CorrectAnswer)
	given := strings.TrimSpace(answer)
	// Open-ended activities (no expected answer) accept any thoughtful response.
	isCorrect := (expected == "" && given != "") || (expected != "" && strings.EqualFold(given, expected))
	now := time.Now()
	durationSec := int64(now.Sub(att.StartedAt).Seconds())

	if !isCorrect {
		_ = s.repo.RecordSignal(models.EngagementSignal{
			ID:          uuid.New().String(),
			StudentID:   studentID,
			AttemptID:   attemptID,
			SignalType:  "errors",
			SignalValue: fmt.Sprintf("failed_attempt_%d: %s", attempts, answer),
		})
		return att, false, nil
	}

	if err := s.repo.UpdateAttempt(attemptID, "completed", now, durationSec, attempts); err != nil {
		return nil, false, err
	}

	_ = s.repo.RecordSignal(models.EngagementSignal{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		AttemptID:   attemptID,
		SignalType:  "completion",
		SignalValue: "success",
	})
	_ = s.repo.RecordSignal(models.EngagementSignal{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		AttemptID:   attemptID,
		SignalType:  "duration",
		SignalValue: fmt.Sprintf("%d seconds", durationSec),
	})
	if attempts > 1 {
		_ = s.repo.RecordSignal(models.EngagementSignal{
			ID:          uuid.New().String(),
			StudentID:   studentID,
			AttemptID:   attemptID,
			SignalType:  "repeated_attempts",
			SignalValue: fmt.Sprintf("%d tries", attempts),
		})
	}

	att.CompletedAt = &now
	att.Duration = durationSec
	att.Status = "completed"
	att.AttemptCount = attempts

	return att, true, nil
}

func (s *Service) AbandonActivity(activityID string, attemptID string, studentID string, attempts int) (*models.ActivityAttempt, error) {
	att, err := s.repo.GetAttempt(attemptID)
	if err != nil {
		return nil, err
	}
	if att.Status != "in_progress" {
		return nil, errors.New("attempt already finished")
	}

	now := time.Now()
	durationSec := int64(now.Sub(att.StartedAt).Seconds())

	if err := s.repo.UpdateAttempt(attemptID, "abandoned", now, durationSec, attempts); err != nil {
		return nil, err
	}

	_ = s.repo.RecordSignal(models.EngagementSignal{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		AttemptID:   attemptID,
		SignalType:  "abandonment",
		SignalValue: "user_quit",
	})
	_ = s.repo.RecordSignal(models.EngagementSignal{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		AttemptID:   attemptID,
		SignalType:  "duration",
		SignalValue: fmt.Sprintf("%d seconds", durationSec),
	})

	att.CompletedAt = &now
	att.Duration = durationSec
	att.Status = "abandoned"
	att.AttemptCount = attempts

	return att, nil
}

func (s *Service) GetStudentHistory(studentID string) ([]models.StudentHistoryItem, error) {
	return s.repo.GetStudentHistory(studentID)
}
