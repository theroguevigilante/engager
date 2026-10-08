package intervention

import (
	"testing"

	"github.com/adaptive-learning/backend/internal/assessment"
	"github.com/adaptive-learning/backend/internal/models"
)

func TestInterventionPolicy(t *testing.T) {
	mockAttempt := &models.ActivityAttempt{ID: "att-1", AttemptCount: 1, Status: "completed"}

	tests := []struct {
		name           string
		assessment     models.EngagementAssessment
		expectedType   models.InterventionType
		expectedReason string
	}{
		{
			name: "ENGAGED -> CONTINUE",
			assessment: models.EngagementAssessment{
				State:      models.StateEngaged,
				Confidence: 0.85,
				Reasons:    []string{assessment.ReasonNormalPace, assessment.ReasonCompletedActivity},
			},
			expectedType:   models.InterventionContinue,
			expectedReason: "state_fallback_engaged",
		},
		{
			name: "BORED + rapid_completion -> INCREASE_CHALLENGE",
			assessment: models.EngagementAssessment{
				State:      models.StateBored,
				Confidence: 0.90,
				Reasons:    []string{assessment.ReasonRapidCompletion, assessment.ReasonHighActiveRatio},
			},
			expectedType:   models.InterventionIncreaseChallenge,
			expectedReason: assessment.ReasonRapidCompletion,
		},
		{
			name: "STRUGGLING + repeated_attempts -> SCAFFOLD",
			assessment: models.EngagementAssessment{
				State:      models.StateStruggling,
				Confidence: 0.90,
				Reasons:    []string{assessment.ReasonRepeatedAttempts, assessment.ReasonExtendedDuration},
			},
			expectedType:   models.InterventionScaffold,
			expectedReason: assessment.ReasonRepeatedAttempts,
		},
		{
			name: "DISENGAGED + high_idle_ratio -> RE_ENGAGE",
			assessment: models.EngagementAssessment{
				State:      models.StateDisengaged,
				Confidence: 0.80,
				Reasons:    []string{assessment.ReasonHighIdleRatio},
			},
			expectedType:   models.InterventionReEngage,
			expectedReason: assessment.ReasonHighIdleRatio,
		},
		{
			name: "DISENGAGED + multiple_tab_switches -> CHANGE_ACTIVITY",
			assessment: models.EngagementAssessment{
				State:      models.StateDisengaged,
				Confidence: 0.85,
				Reasons:    []string{assessment.ReasonMultipleTabSwitches},
			},
			expectedType:   models.InterventionChangeActivity,
			expectedReason: assessment.ReasonMultipleTabSwitches,
		},
		{
			name: "Multiple reasons: repeated_attempts beats high_idle_ratio and rapid_completion",
			assessment: models.EngagementAssessment{
				State:      models.StateStruggling,
				Confidence: 0.95,
				Reasons:    []string{assessment.ReasonRapidCompletion, assessment.ReasonRepeatedAttempts, assessment.ReasonHighIdleRatio},
			},
			expectedType:   models.InterventionScaffold,
			expectedReason: assessment.ReasonRepeatedAttempts,
		},
		{
			name: "Unknown / unexpected reason -> safe state fallback",
			assessment: models.EngagementAssessment{
				State:      models.StateStruggling,
				Confidence: 0.65,
				Reasons:    []string{"unmapped_experimental_signal"},
			},
			expectedType:   models.InterventionScaffold,
			expectedReason: "state_fallback_struggling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iType, reason, _ := SelectIntervention(&tt.assessment, mockAttempt)
			if iType != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, iType)
			}
			if reason != tt.expectedReason {
				t.Errorf("expected reason %s, got %s", tt.expectedReason, reason)
			}
		})
	}
}
