package models

import "time"

type LoopEffectiveness string

const (
	EffectivenessEffective   LoopEffectiveness = "EFFECTIVE"
	EffectivenessIneffective LoopEffectiveness = "INEFFECTIVE"
	EffectivenessNeutral     LoopEffectiveness = "NEUTRAL"
)

type AdaptationLoop struct {
	ID                     string            `json:"id"`
	StudentID              string            `json:"student_id"`
	OriginalAttemptID      string            `json:"original_attempt_id"`
	AdaptedAttemptID       string            `json:"adapted_attempt_id"`
	AssessmentID           string            `json:"assessment_id"`
	InterventionID         string            `json:"intervention_id"`
	PersonalizedActivityID string            `json:"personalized_activity_id"`
	PreState               EngagementState   `json:"pre_state"`
	PostState              EngagementState   `json:"post_state"`
	PreEngagementScore     float64           `json:"pre_engagement_score"`
	PostEngagementScore    float64           `json:"post_engagement_score"`
	EngagementDelta        float64           `json:"engagement_delta"`
	Effectiveness          LoopEffectiveness `json:"effectiveness"`
	EffectivenessReasons   []string          `json:"effectiveness_reasons"`
	CreatedAt              time.Time         `json:"created_at"`
}
