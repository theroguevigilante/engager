package intervention

import (
	"github.com/adaptive-learning/backend/internal/assessment"
	"github.com/adaptive-learning/backend/internal/models"
)

type PriorityRule struct {
	ReasonCode       string
	InterventionType models.InterventionType
}

// Deterministic priority ordering for reason codes
var PriorityRules = []PriorityRule{
	{ReasonCode: assessment.ReasonRepeatedAttempts, InterventionType: models.InterventionScaffold},
	{ReasonCode: assessment.ReasonExtendedDuration, InterventionType: models.InterventionScaffold},
	{ReasonCode: assessment.ReasonHighIdleRatio, InterventionType: models.InterventionReEngage},
	{ReasonCode: assessment.ReasonActivityAbandoned, InterventionType: models.InterventionChangeActivity},
	{ReasonCode: assessment.ReasonMultipleTabSwitches, InterventionType: models.InterventionChangeActivity},
	{ReasonCode: assessment.ReasonRapidCompletion, InterventionType: models.InterventionIncreaseChallenge},
	{ReasonCode: assessment.ReasonHighActiveRatio, InterventionType: models.InterventionIncreaseChallenge},
}

// SelectIntervention runs the reason priority evaluation with state fallback
func SelectIntervention(ass *models.EngagementAssessment, att *models.ActivityAttempt) (models.InterventionType, string, float64) {
	reasonMap := make(map[string]bool)
	for _, r := range ass.Reasons {
		reasonMap[r] = true
	}

	// 1. Evaluate explicit priority rules based on assessment reasons
	for _, rule := range PriorityRules {
		if reasonMap[rule.ReasonCode] {
			return rule.InterventionType, rule.ReasonCode, ass.Confidence
		}
	}

	// 2. Safe deterministic fallback based on engagement state
	switch ass.State {
	case models.StateBored:
		return models.InterventionIncreaseChallenge, "state_fallback_bored", ass.Confidence
	case models.StateStruggling:
		return models.InterventionScaffold, "state_fallback_struggling", ass.Confidence
	case models.StateDisengaged:
		return models.InterventionReEngage, "state_fallback_disengaged", ass.Confidence
	case models.StateEngaged:
		fallthrough
	default:
		return models.InterventionContinue, "state_fallback_engaged", ass.Confidence
	}
}
