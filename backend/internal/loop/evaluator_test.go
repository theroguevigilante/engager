package loop

import (
	"testing"

	"github.com/adaptive-learning/backend/internal/models"
)

func TestLoopEvaluationScenarios(t *testing.T) {
	tests := []struct {
		name                  string
		preAssessment         models.EngagementAssessment
		postAssessment        models.EngagementAssessment
		postAttempt           models.ActivityAttempt
		interventionType      models.InterventionType
		expectedEffectiveness models.LoopEffectiveness
	}{
		{
			name: "Scenario A: BORED -> INCREASE_CHALLENGE -> ENGAGED (Completed) => EFFECTIVE",
			preAssessment: models.EngagementAssessment{
				State:   models.StateBored,
				Metrics: models.AssessmentMetrics{SessionDuration: 10, ActiveRatio: 0.90, IdleRatio: 0.10, InteractionRate: 0.20},
			},
			postAssessment: models.EngagementAssessment{
				State:   models.StateEngaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 45, ActiveRatio: 0.85, IdleRatio: 0.15, InteractionRate: 0.60},
			},
			postAttempt:           models.ActivityAttempt{Status: "completed"},
			interventionType:      models.InterventionIncreaseChallenge,
			expectedEffectiveness: models.EffectivenessEffective,
		},
		{
			name: "Scenario B: STRUGGLING -> SCAFFOLD -> Still STRUGGLING (Completed) => NEUTRAL",
			preAssessment: models.EngagementAssessment{
				State:   models.StateStruggling,
				Metrics: models.AssessmentMetrics{SessionDuration: 60, ActiveRatio: 0.70, IdleRatio: 0.30, InteractionRate: 0.90},
			},
			postAssessment: models.EngagementAssessment{
				State:   models.StateStruggling,
				Metrics: models.AssessmentMetrics{SessionDuration: 55, ActiveRatio: 0.72, IdleRatio: 0.28, InteractionRate: 0.85},
			},
			postAttempt:           models.ActivityAttempt{Status: "completed"},
			interventionType:      models.InterventionScaffold,
			expectedEffectiveness: models.EffectivenessNeutral,
		},
		{
			name: "Scenario C: DISENGAGED -> RE_ENGAGE -> DISENGAGED (Abandoned) => INEFFECTIVE",
			preAssessment: models.EngagementAssessment{
				State:   models.StateDisengaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 40, ActiveRatio: 0.20, IdleRatio: 0.80, InteractionRate: 0.05},
			},
			postAssessment: models.EngagementAssessment{
				State:   models.StateDisengaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 25, ActiveRatio: 0.10, IdleRatio: 0.90, InteractionRate: 0.02},
			},
			postAttempt:           models.ActivityAttempt{Status: "abandoned"},
			interventionType:      models.InterventionReEngage,
			expectedEffectiveness: models.EffectivenessIneffective,
		},
		{
			name: "Scenario D: ENGAGED -> CONTINUE -> ENGAGED (Completed) => EFFECTIVE",
			preAssessment: models.EngagementAssessment{
				State:   models.StateEngaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 35, ActiveRatio: 0.80, IdleRatio: 0.20, InteractionRate: 0.40},
			},
			postAssessment: models.EngagementAssessment{
				State:   models.StateEngaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 40, ActiveRatio: 0.82, IdleRatio: 0.18, InteractionRate: 0.45},
			},
			postAttempt:           models.ActivityAttempt{Status: "completed"},
			interventionType:      models.InterventionContinue,
			expectedEffectiveness: models.EffectivenessEffective,
		},
		{
			name: "Scenario E: Insufficient Post Telemetry => NEUTRAL",
			preAssessment: models.EngagementAssessment{
				State:   models.StateBored,
				Metrics: models.AssessmentMetrics{SessionDuration: 15, ActiveRatio: 0.85, IdleRatio: 0.15, InteractionRate: 0.30},
			},
			postAssessment: models.EngagementAssessment{
				State:   models.StateEngaged,
				Metrics: models.AssessmentMetrics{SessionDuration: 2, ActiveRatio: 0.0, IdleRatio: 0.0, InteractionRate: 0.0},
			},
			postAttempt:           models.ActivityAttempt{Status: "completed"},
			interventionType:      models.InterventionIncreaseChallenge,
			expectedEffectiveness: models.EffectivenessNeutral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eff, _, _, _, reasons := EvaluateLoop(&tt.preAssessment, &tt.postAssessment, &tt.postAttempt, tt.interventionType)
			if eff != tt.expectedEffectiveness {
				t.Errorf("expected effectiveness %s, got %s (reasons: %v)", tt.expectedEffectiveness, eff, reasons)
			}
			if len(reasons) == 0 {
				t.Errorf("expected explicit reasons to be provided")
			}
		})
	}
}
