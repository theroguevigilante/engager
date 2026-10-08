package models

// TelemetryBatch matches the incoming JSON payload from the frontend tracker
type TelemetryBatch struct {
	StudentID        string `json:"student_id"`
	SessionDuration  int    `json:"session_duration"`
	ActiveDuration   int    `json:"active_duration"`
	IdleDuration     int    `json:"idle_duration"`
	TabSwitches      int    `json:"tab_switches"`
	InteractionCount int    `json:"interaction_count"`
}
