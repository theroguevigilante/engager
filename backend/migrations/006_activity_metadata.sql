-- Author-facing activity metadata used by the Activity Library
ALTER TABLE activities ADD COLUMN topic TEXT NOT NULL DEFAULT 'General';
ALTER TABLE activities ADD COLUMN difficulty TEXT NOT NULL DEFAULT 'intermediate';
ALTER TABLE activities ADD COLUMN estimated_duration_minutes INTEGER NOT NULL DEFAULT 5;
