package models

import "time"

type PersonalizedActivity struct {
	ID                    string           `json:"id"`
	AttemptID             string           `json:"attempt_id"`
	SourceActivityID      string           `json:"source_activity_id"`
	InterventionID        string           `json:"intervention_id"`
	Title                 string           `json:"title"`
	Content               string           `json:"content"`
	Difficulty            string           `json:"difficulty"`
	InterventionType      InterventionType `json:"intervention_type"`
	PersonalizationReason string           `json:"personalization_reason"`
	IsGenerated           bool             `json:"is_generated"`
	CreatedAt             time.Time        `json:"created_at"`
}
