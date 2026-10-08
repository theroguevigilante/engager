package assessment

import (
	"math"

	"github.com/adaptive-learning/backend/internal/models"
)

// Heuristic configuration thresholds
const (
	MinTelemetryDurationSec = 5

	// Boredom thresholds
	BoredomMaxDurationSec     = 20
	BoredomMinActiveRatio     = 0.70
	BoredomMaxInteractionRate = 0.50

	// Struggling thresholds
	StrugglingMinDurationSec     = 45
	StrugglingMinActiveRatio     = 0.60
	StrugglingMinInteractionRate = 0.80
	StrugglingMinAttemptCount    = 2

	// Disengagement thresholds
	DisengagedMinIdleRatio       = 0.50
	DisengagedMaxInteractionRate = 0.15
	DisengagedMinTabSwitches     = 2

	// Engaged thresholds
	EngagedMinActiveRatio     = 0.60
	EngagedMaxTabSwitches     = 2
	EngagedMinInteractionRate = 0.20
)

// Explicit behavioral reason codes
const (
	ReasonRapidCompletion     = "rapid_completion"
	ReasonMinimalInteractions = "minimal_interactions"
	ReasonHighActiveRatio     = "high_active_ratio"
	ReasonExtendedDuration    = "extended_duration"
	ReasonHighInteractionRate = "high_interaction_rate"
	ReasonRepeatedAttempts    = "repeated_attempts"
	ReasonUncompleted         = "uncompleted_activity"
	ReasonHighIdleRatio       = "high_idle_ratio"
	ReasonMultipleTabSwitches = "multiple_tab_switches"
	ReasonLowInteractionRate  = "low_interaction_rate"
	ReasonActivityAbandoned   = "activity_abandoned"
	ReasonNormalPace          = "normal_active_pace"
	ReasonCompletedActivity   = "completed_activity"
	ReasonInsufficientData    = "insufficient_telemetry_data"
)

func ComputeMetrics(signals []models.EngagementSignal, attempt models.ActivityAttempt) models.AssessmentMetrics {
	var sessionDuration, activeDuration, idleDuration, tabSwitches, interactionCount int

	for _, sig := range signals {
		val := parseInt(sig.SignalValue)
		switch sig.SignalType {
		case "session_duration":
			if val > sessionDuration {
				sessionDuration = val
			}
		case "active_duration":
			if val > activeDuration {
				activeDuration = val
			}
		case "idle_duration":
			if val > idleDuration {
				idleDuration = val
			}
		case "tab_switches":
			if val > tabSwitches {
				tabSwitches = val
			}
		case "interaction_count":
			if val > interactionCount {
				interactionCount = val
			}
		}
	}

	if sessionDuration == 0 && attempt.Duration > 0 {
		sessionDuration = int(attempt.Duration)
	}

	var activeRatio, idleRatio, interactionRate, tabSwitchRate float64
	if sessionDuration > 0 {
		activeRatio = float64(activeDuration) / float64(sessionDuration)
		idleRatio = float64(idleDuration) / float64(sessionDuration)
		interactionRate = float64(interactionCount) / float64(sessionDuration)
		tabSwitchRate = float64(tabSwitches) / float64(sessionDuration)
	}

	return models.AssessmentMetrics{
		SessionDuration:  sessionDuration,
		ActiveDuration:   activeDuration,
		IdleDuration:     idleDuration,
		TabSwitches:      tabSwitches,
		InteractionCount: interactionCount,
		AttemptCount:     attempt.AttemptCount,
		ActiveRatio:      roundFloat(activeRatio, 2),
		IdleRatio:        roundFloat(idleRatio, 2),
		InteractionRate:  roundFloat(interactionRate, 2),
		TabSwitchRate:    roundFloat(tabSwitchRate, 2),
		Status:           attempt.Status,
	}
}

func Classify(m models.AssessmentMetrics) (models.EngagementState, float64, []string) {
	if m.SessionDuration < MinTelemetryDurationSec && m.InteractionCount == 0 {
		return models.StateEngaged, 0.20, []string{ReasonInsufficientData}
	}

	var reasons []string

	// Rule 1: DISENGAGED
	isAbandoned := m.Status == "abandoned"
	hasHighIdle := m.IdleRatio >= DisengagedMinIdleRatio
	hasTabSwitches := m.TabSwitches >= DisengagedMinTabSwitches
	hasLowInteractions := m.InteractionRate <= DisengagedMaxInteractionRate

	if isAbandoned || (hasHighIdle && (hasTabSwitches || hasLowInteractions)) {
		if isAbandoned {
			reasons = append(reasons, ReasonActivityAbandoned)
		}
		if hasHighIdle {
			reasons = append(reasons, ReasonHighIdleRatio)
		}
		if hasTabSwitches {
			reasons = append(reasons, ReasonMultipleTabSwitches)
		}
		if hasLowInteractions {
			reasons = append(reasons, ReasonLowInteractionRate)
		}
		conf := 0.70
		if isAbandoned && hasHighIdle {
			conf = 0.95
		} else if len(reasons) >= 3 {
			conf = 0.85
		}
		return models.StateDisengaged, conf, reasons
	}

	// Rule 2: STRUGGLING
	hasLongDuration := m.SessionDuration >= StrugglingMinDurationSec
	hasHighInteraction := m.InteractionRate >= StrugglingMinInteractionRate
	hasMultiAttempts := m.AttemptCount >= StrugglingMinAttemptCount
	isUnfinished := m.Status != "completed"

	if (hasLongDuration && hasHighInteraction && isUnfinished) || (hasMultiAttempts && hasLongDuration) {
		if hasLongDuration {
			reasons = append(reasons, ReasonExtendedDuration)
		}
		if hasHighInteraction {
			reasons = append(reasons, ReasonHighInteractionRate)
		}
		if hasMultiAttempts {
			reasons = append(reasons, ReasonRepeatedAttempts)
		}
		if isUnfinished {
			reasons = append(reasons, ReasonUncompleted)
		}
		conf := 0.75
		if hasMultiAttempts && isUnfinished {
			conf = 0.90
		}
		return models.StateStruggling, conf, reasons
	}

	// Rule 3: BORED / UNDER-CHALLENGED
	isRapid := m.SessionDuration <= BoredomMaxDurationSec && m.Status == "completed"
	hasHighActive := m.ActiveRatio >= BoredomMinActiveRatio
	hasLowInteractionRate := m.InteractionRate <= BoredomMaxInteractionRate

	if isRapid && hasHighActive && hasLowInteractionRate {
		reasons = append(reasons, ReasonRapidCompletion, ReasonHighActiveRatio, ReasonMinimalInteractions)
		conf := 0.85
		if m.AttemptCount <= 1 {
			conf = 0.90
		}
		return models.StateBored, conf, reasons
	}

	// Rule 4: ENGAGED (Default baseline)
	if m.ActiveRatio >= EngagedMinActiveRatio {
		reasons = append(reasons, ReasonHighActiveRatio)
	}
	if m.Status == "completed" {
		reasons = append(reasons, ReasonCompletedActivity)
	}
	reasons = append(reasons, ReasonNormalPace)

	conf := 0.80
	if m.ActiveRatio < EngagedMinActiveRatio {
		conf = 0.50
	}
	return models.StateEngaged, conf, reasons
}

func parseInt(s string) int {
	var val int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
		} else {
			break
		}
	}
	return val
}

func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
