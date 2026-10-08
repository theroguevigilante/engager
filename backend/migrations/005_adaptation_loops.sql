-- Link adapted attempts back to their originating attempt
ALTER TABLE activity_attempts ADD COLUMN parent_attempt_id TEXT REFERENCES activity_attempts(id);

CREATE TABLE IF NOT EXISTS adaptation_loops (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id),
    original_attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    adapted_attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    assessment_id TEXT NOT NULL REFERENCES engagement_assessments(id),
    intervention_id TEXT NOT NULL REFERENCES attempt_interventions(id),
    personalized_activity_id TEXT NOT NULL REFERENCES personalized_activities(id),
    pre_state TEXT NOT NULL,
    post_state TEXT NOT NULL,
    pre_engagement_score REAL NOT NULL,
    post_engagement_score REAL NOT NULL,
    engagement_delta REAL NOT NULL,
    effectiveness TEXT NOT NULL, -- 'EFFECTIVE', 'INEFFECTIVE', 'NEUTRAL'
    effectiveness_reasons_json TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_adaptation_loops_original_attempt ON adaptation_loops(original_attempt_id);
CREATE INDEX IF NOT EXISTS idx_adaptation_loops_student ON adaptation_loops(student_id);
