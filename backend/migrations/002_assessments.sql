CREATE TABLE IF NOT EXISTS engagement_assessments (
    id TEXT PRIMARY KEY,
    attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    state TEXT NOT NULL, -- 'ENGAGED', 'STRUGGLING', 'BORED', 'DISENGAGED'
    confidence REAL NOT NULL,
    metrics_json TEXT NOT NULL,
    reasons_json TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_engagement_assessments_attempt_id ON engagement_assessments(attempt_id);
