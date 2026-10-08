package models

import "time"

type Student struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Activity struct {
	ID                       string    `json:"id"`
	Title                    string    `json:"title"`
	Description              string    `json:"description"`
	Content                  string    `json:"content"`
	CorrectAnswer            string    `json:"correct_answer,omitempty"`
	Topic                    string    `json:"topic,omitempty"`
	Difficulty               string    `json:"difficulty,omitempty"`
	EstimatedDurationMinutes int       `json:"estimated_duration_minutes,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
}

type ActivityAttempt struct {
	ID              string     `json:"id"`
	StudentID       string     `json:"student_id"`
	ActivityID      string     `json:"activity_id"`
	ParentAttemptID *string    `json:"parent_attempt_id,omitempty"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Duration        int64      `json:"duration"`
	Status          string     `json:"status"`
	AttemptCount    int        `json:"attempt_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

type EngagementSignal struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	AttemptID   string    `json:"attempt_id"`
	SignalType  string    `json:"signal_type"`
	SignalValue string    `json:"signal_value"`
	CreatedAt   time.Time `json:"created_at"`
}

type StudentHistoryItem struct {
	AttemptID     string     `json:"attempt_id"`
	ActivityID    string     `json:"activity_id"`
	ActivityTitle string     `json:"activity_title"`
	Status        string     `json:"status"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	Duration      int64      `json:"duration"`
	AttemptCount  int        `json:"attempt_count"`
	SignalsCount  int        `json:"signals_count"`
}
