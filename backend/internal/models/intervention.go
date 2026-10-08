package models

import "time"

type InterventionType string

const (
	InterventionContinue          InterventionType = "CONTINUE"
	InterventionIncreaseChallenge InterventionType = "INCREASE_CHALLENGE"
	InterventionScaffold          InterventionType = "SCAFFOLD"
	InterventionReEngage          InterventionType = "RE_ENGAGE"
	InterventionChangeActivity    InterventionType = "CHANGE_ACTIVITY"
)

type InterventionMetadata struct {
	Type        InterventionType `json:"type"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Category    string           `json:"category"`
}

type AttemptIntervention struct {
	ID               string               `json:"id"`
	AttemptID        string               `json:"attempt_id"`
	AssessmentID     string               `json:"assessment_id"`
	InterventionType InterventionType     `json:"intervention_type"`
	Reason           string               `json:"reason"`
	Confidence       float64              `json:"confidence"`
	Metadata         InterventionMetadata `json:"metadata"`
	CreatedAt        time.Time            `json:"created_at"`
}
