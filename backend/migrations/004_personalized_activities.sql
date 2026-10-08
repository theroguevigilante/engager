CREATE TABLE IF NOT EXISTS personalized_activities (
    id TEXT PRIMARY KEY,
    attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    source_activity_id TEXT NOT NULL REFERENCES activities(id),
    intervention_id TEXT NOT NULL REFERENCES attempt_interventions(id),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    personalization_reason TEXT NOT NULL,
    is_generated BOOLEAN NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_personalized_activities_attempt_id ON personalized_activities(attempt_id);
