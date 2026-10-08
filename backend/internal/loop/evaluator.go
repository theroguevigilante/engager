package loop

import (
	"math"

	"github.com/adaptive-learning/backend/internal/models"
)

// Named constants for engagement scoring and delta thresholds
const (
	MinMeaningfulDelta     = 0.15
	SubstantialWorsenDelta = -0.15
	MinTelemetryDuration   = 5

	ReasonStateImprovedToEngaged       = "state_improved_to_engaged"
	ReasonHigherActiveRatio            = "higher_active_ratio"
	ReasonLowerIdleRatio               = "lower_idle_ratio"
	ReasonActivityCompletedAfterInterv = "activity_completed_after_intervention"
	ReasonActivityAbandoned            = "activity_abandoned"
	ReasonContinuedDisengagement       = "continued_disengagement"
	ReasonEngagementUnchanged          = "engagement_unchanged"
	ReasonInsufficientPostIntervData   = "insufficient_post_intervention_data"
	ReasonMaintainedHealthyEngagement  = "maintained_healthy_engagement"
)

// CalculateEngagementScore computes a deterministic 0.0 - 1.0 composite index
func CalculateEngagementScore(m models.AssessmentMetrics, state models.EngagementState) float64 {
	if m.SessionDuration < MinTelemetryDuration && m.InteractionCount == 0 {
		return 0.50
	}

	score := (m.ActiveRatio * 0.50) + ((1.0 - m.IdleRatio) * 0.30)

	// Interaction bonus up to 0.20
	intRate := m.InteractionRate
	if intRate > 1.0 {
		intRate = 1.0
	}
	score += (intRate * 0.20)

	// Tab penalty up to -0.15
	if m.TabSwitches > 0 {
		penalty := float64(m.TabSwitches) * 0.05
		if penalty > 0.15 {
			penalty = 0.15
		}
		score -= penalty
	}

	// State modifiers
	switch state {
	case models.StateEngaged:
		score += 0.10
	case models.StateBored:
		score -= 0.05
	case models.StateStruggling:
		score -= 0.15
	case models.StateDisengaged:
		score -= 0.25
	}

	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return math.Round(score*100) / 100
}

// EvaluateLoop evaluates the transition and produces a deterministic verdict
func EvaluateLoop(
	preAss *models.EngagementAssessment,
	postAss *models.EngagementAssessment,
	postAtt *models.ActivityAttempt,
	interType models.InterventionType,
) (models.LoopEffectiveness, float64, float64, float64, []string) {
	var reasons []string

	if postAss.Metrics.SessionDuration < MinTelemetryDuration && postAss.Metrics.InteractionCount == 0 {
		return models.EffectivenessNeutral, 0.50, 0.50, 0.0, []string{ReasonInsufficientPostIntervData}
	}

	preScore := CalculateEngagementScore(preAss.Metrics, preAss.State)
	postScore := CalculateEngagementScore(postAss.Metrics, postAss.State)
	delta := math.Round((postScore-preScore)*100) / 100

	isAbandoned := postAtt.Status == "abandoned"
	isCompleted := postAtt.Status == "completed"

	if postAss.Metrics.ActiveRatio > preAss.Metrics.ActiveRatio+0.10 {
		reasons = append(reasons, ReasonHigherActiveRatio)
	}
	if postAss.Metrics.IdleRatio < preAss.Metrics.IdleRatio-0.10 {
		reasons = append(reasons, ReasonLowerIdleRatio)
	}

	// INEFFECTIVE rules
	if isAbandoned {
		reasons = append(reasons, ReasonActivityAbandoned)
		return models.EffectivenessIneffective, preScore, postScore, delta, reasons
	}

	if postAss.State == models.StateDisengaged && (preAss.State == models.StateDisengaged || preAss.State == models.StateStruggling) {
		reasons = append(reasons, ReasonContinuedDisengagement)
		return models.EffectivenessIneffective, preScore, postScore, delta, reasons
	}

	if delta <= SubstantialWorsenDelta {
		return models.EffectivenessIneffective, preScore, postScore, delta, reasons
	}

	// EFFECTIVE rules
	if isCompleted {
		if preAss.State != models.StateEngaged && postAss.State == models.StateEngaged {
			reasons = append(reasons, ReasonStateImprovedToEngaged, ReasonActivityCompletedAfterInterv)
			return models.EffectivenessEffective, preScore, postScore, delta, reasons
		}

		if delta >= MinMeaningfulDelta {
			reasons = append(reasons, ReasonActivityCompletedAfterInterv)
			return models.EffectivenessEffective, preScore, postScore, delta, reasons
		}

		if interType == models.InterventionContinue && preAss.State == models.StateEngaged && postAss.State == models.StateEngaged {
			reasons = append(reasons, ReasonMaintainedHealthyEngagement, ReasonActivityCompletedAfterInterv)
			return models.EffectivenessEffective, preScore, postScore, delta, reasons
		}
	}

	// Default NEUTRAL
	if len(reasons) == 0 {
		reasons = append(reasons, ReasonEngagementUnchanged)
	}
	return models.EffectivenessNeutral, preScore, postScore, delta, reasons
}
