package assessment

import (
	"testing"

	"github.com/adaptive-learning/backend/internal/models"
)

func TestClassifierScenarios(t *testing.T) {
	tests := []struct {
		name          string
		metrics       models.AssessmentMetrics
		expectedState models.EngagementState
		minConfidence float64
	}{
		{
			name: "Scenario A: Bored / Under-challenged",
			metrics: models.AssessmentMetrics{
				SessionDuration:  12,
				ActiveDuration:   11,
				IdleDuration:     1,
				TabSwitches:      0,
				InteractionCount: 2,
				AttemptCount:     1,
				ActiveRatio:      0.92,
				IdleRatio:        0.08,
				InteractionRate:  0.17,
				Status:           "completed",
			},
			expectedState: models.StateBored,
			minConfidence: 0.80,
		},
		{
			name: "Scenario B: Struggling",
			metrics: models.AssessmentMetrics{
				SessionDuration:  70,
				ActiveDuration:   55,
				IdleDuration:     15,
				TabSwitches:      1,
				InteractionCount: 65,
				AttemptCount:     2,
				ActiveRatio:      0.79,
				IdleRatio:        0.21,
				InteractionRate:  0.93,
				Status:           "in_progress",
			},
			expectedState: models.StateStruggling,
			minConfidence: 0.70,
		},
		{
			name: "Scenario C: Disengaged",
			metrics: models.AssessmentMetrics{
				SessionDuration:  60,
				ActiveDuration:   15,
				IdleDuration:     45,
				TabSwitches:      4,
				InteractionCount: 3,
				AttemptCount:     1,
				ActiveRatio:      0.25,
				IdleRatio:        0.75,
				InteractionRate:  0.05,
				Status:           "abandoned",
			},
			expectedState: models.StateDisengaged,
			minConfidence: 0.85,
		},
		{
			name: "Scenario D: Engaged",
			metrics: models.AssessmentMetrics{
				SessionDuration:  40,
				ActiveDuration:   32,
				IdleDuration:     8,
				TabSwitches:      0,
				InteractionCount: 16,
				AttemptCount:     1,
				ActiveRatio:      0.80,
				IdleRatio:        0.20,
				InteractionRate:  0.40,
				Status:           "completed",
			},
			expectedState: models.StateEngaged,
			minConfidence: 0.75,
		},
		{
			name: "Scenario E: Minimal / Missing Telemetry",
			metrics: models.AssessmentMetrics{
				SessionDuration:  2,
				ActiveDuration:   0,
				IdleDuration:     0,
				TabSwitches:      0,
				InteractionCount: 0,
				AttemptCount:     1,
				ActiveRatio:      0.0,
				IdleRatio:        0.0,
				InteractionRate:  0.0,
				Status:           "in_progress",
			},
			expectedState: models.StateEngaged,
			minConfidence: 0.20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, conf, reasons := Classify(tt.metrics)
			if state != tt.expectedState {
				t.Errorf("expected state %s, got %s (reasons: %v)", tt.expectedState, state, reasons)
			}
			if conf < tt.minConfidence {
				t.Errorf("expected confidence >= %f, got %f", tt.minConfidence, conf)
			}
			if len(reasons) == 0 {
				t.Errorf("expected explicit reasons to be provided")
			}
		})
	}
}
