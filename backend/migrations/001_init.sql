CREATE TABLE IF NOT EXISTS students (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activities (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activity_attempts (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id),
    activity_id TEXT NOT NULL REFERENCES activities(id),
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    duration INTEGER DEFAULT 0,
    status TEXT NOT NULL, -- 'in_progress', 'completed', 'abandoned'
    attempt_count INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS engagement_signals (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id),
    attempt_id TEXT NOT NULL REFERENCES activity_attempts(id),
    signal_type TEXT NOT NULL, -- 'completion', 'abandonment', 'repeated_attempts', 'errors', 'duration'
    signal_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
