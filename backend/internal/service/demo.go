package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adaptive-learning/backend/internal/demo"
	"github.com/adaptive-learning/backend/internal/models"
	"github.com/adaptive-learning/backend/internal/planner"
)

// DemoResult summarises one scripted end-to-end adaptation loop.
type DemoResult struct {
	ActivityTitle      string                 `json:"activity_title"`
	Scenario           string                 `json:"scenario"`
	Need               string                 `json:"need"`
	PreState           models.EngagementState `json:"pre_state"`
	PreReasons         []string               `json:"pre_reasons"`
	InterventionType   string                 `json:"intervention_type"`
	InterventionReason string                 `json:"intervention_reason"`
	AdaptedTitle       string                 `json:"adapted_title"`
	AdaptedBy          string                 `json:"adapted_by"`
	AdaptedContent     string                 `json:"adapted_content"`
	PostState          models.EngagementState `json:"post_state"`
	Loop               *models.AdaptationLoop `json:"loop"`
}

func (s *Service) playTelemetry(ctx context.Context, activityID, studentID string, parent *string, tel demo.Telemetry) (*models.ActivityAttempt, error) {
	att, err := s.StartActivity(activityID, studentID, parent)
	if err != nil {
		return nil, err
	}
	if err := s.IngestSignals(ctx, att.ID, models.TelemetryBatch{
		StudentID: studentID, SessionDuration: tel.Session, ActiveDuration: tel.Active,
		IdleDuration: tel.Idle, TabSwitches: tel.TabSwitches, InteractionCount: tel.Interactions,
	}); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAttempt(att.ID, tel.Status, time.Now(), int64(tel.Session), tel.Attempts); err != nil {
		return nil, err
	}
	return att, nil
}

// RunDemoLoop drives the full closed loop with scripted behavioural telemetry
// for a chosen scenario (struggling / bored / drifting) and an optional
// self-reported need. All assessment, intervention and outcome decisions still
// come from the real deterministic engine.
func (s *Service) RunDemoLoop(ctx context.Context, activityID, studentID, scenario, need string) (*DemoResult, error) {
	act, err := s.repo.GetActivity(activityID)
	if err != nil {
		return nil, fmt.Errorf("activity not found: %w", err)
	}
	sc := demo.Find(scenario)

	orig, err := s.playTelemetry(ctx, activityID, studentID, nil, sc.Pre)
	if err != nil {
		return nil, err
	}
	pre, err := s.EvaluateAttempt(ctx, orig.ID)
	if err != nil {
		return nil, err
	}
	interv, err := s.SelectAndSaveIntervention(ctx, orig.ID)
	if err != nil {
		return nil, err
	}
	if need != "" {
		if err := s.RecordCheckin(studentID, need); err != nil {
			return nil, err
		}
	}
	pers, err := s.PersonalizeActivityForAttempt(ctx, orig.ID)
	if err != nil {
		return nil, err
	}

	adapted, err := s.playTelemetry(ctx, activityID, studentID, &orig.ID, sc.Post)
	if err != nil {
		return nil, err
	}
	post, err := s.EvaluateAttempt(ctx, adapted.ID)
	if err != nil {
		return nil, err
	}
	lp, err := s.CloseAdaptationLoop(ctx, orig.ID)
	if err != nil {
		return nil, err
	}

	return &DemoResult{
		ActivityTitle:      act.Title,
		Scenario:           sc.Key,
		Need:               need,
		PreState:           pre.State,
		PreReasons:         pre.Reasons,
		InterventionType:   string(interv.InterventionType),
		InterventionReason: interv.Reason,
		AdaptedTitle:       pers.Title,
		AdaptedBy:          planner.GetLLMStatus().LastSource,
		AdaptedContent:     pers.Content,
		PostState:          post.State,
		Loop:               lp,
	}, nil
}
