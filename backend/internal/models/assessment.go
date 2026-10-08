package models

import "time"

type EngagementState string

const (
	StateEngaged    EngagementState = "ENGAGED"
	StateStruggling EngagementState = "STRUGGLING"
	StateBored      EngagementState = "BORED"
	StateDisengaged EngagementState = "DISENGAGED"
)

type AssessmentMetrics struct {
	SessionDuration  int     `json:"session_duration"`
	ActiveDuration   int     `json:"active_duration"`
	IdleDuration     int     `json:"idle_duration"`
	TabSwitches      int     `json:"tab_switches"`
	InteractionCount int     `json:"interaction_count"`
	AttemptCount     int     `json:"attempt_count"`
	ActiveRatio      float64 `json:"active_ratio"`
	IdleRatio        float64 `json:"idle_ratio"`
	InteractionRate  float64 `json:"interaction_rate"`
	TabSwitchRate    float64 `json:"tab_switch_rate"`
	Status           string  `json:"status"`
}

type EngagementAssessment struct {
	ID         string            `json:"id"`
	AttemptID  string            `json:"attempt_id"`
	State      EngagementState   `json:"state"`
	Confidence float64           `json:"confidence"`
	Metrics    AssessmentMetrics `json:"metrics"`
	Reasons    []string          `json:"reasons"`
	AssessedAt time.Time         `json:"assessed_at"`
}
