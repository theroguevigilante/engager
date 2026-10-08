CREATE TABLE IF NOT EXISTS attempt_interventions (
    id TEXT PRIMARY KEY,
    attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    assessment_id TEXT NOT NULL REFERENCES engagement_assessments(id),
    intervention_type TEXT NOT NULL, -- 'CONTINUE', 'INCREASE_CHALLENGE', 'SCAFFOLD', 'RE_ENGAGE', 'CHANGE_ACTIVITY'
    reason TEXT NOT NULL,
    confidence REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_attempt_interventions_attempt_id ON attempt_interventions(attempt_id);
