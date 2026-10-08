package repository

import (
	"database/sql"
	"time"

	"github.com/adaptive-learning/backend/internal/models"
	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) InitSchema(schemaSQL string) error {
	_, err := r.db.Exec(schemaSQL)
	return err
}

const activityColumns = "id, title, description, content, correct_answer, topic, difficulty, estimated_duration_minutes, created_at"

func scanActivity(sc interface{ Scan(...any) error }) (models.Activity, error) {
	var a models.Activity
	err := sc.Scan(&a.ID, &a.Title, &a.Description, &a.Content, &a.CorrectAnswer,
		&a.Topic, &a.Difficulty, &a.EstimatedDurationMinutes, &a.CreatedAt)
	return a, err
}

// ListActivities returns the library. The expected answer is deliberately
// not exposed to the student-facing listing.
func (r *Repository) ListActivities() ([]models.Activity, error) {
	rows, err := r.db.Query("SELECT " + activityColumns + " FROM activities ORDER BY created_at DESC, title")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Activity
	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return nil, err
		}
		a.CorrectAnswer = ""
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *Repository) GetActivity(id string) (*models.Activity, error) {
	a, err := scanActivity(r.db.QueryRow("SELECT "+activityColumns+" FROM activities WHERE id = ?", id))
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) CreateActivity(a models.Activity) error {
	_, err := r.db.Exec(
		"INSERT INTO activities (id, title, description, content, correct_answer, topic, difficulty, estimated_duration_minutes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		a.ID, a.Title, a.Description, a.Content, a.CorrectAnswer, a.Topic, a.Difficulty, a.EstimatedDurationMinutes,
	)
	return err
}

func (r *Repository) EnsureStudent(id, name string) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO students (id, name) VALUES (?, ?)", id, name)
	return err
}

func (r *Repository) CreateAttempt(att models.ActivityAttempt) error {
	_, err := r.db.Exec(
		"INSERT INTO activity_attempts (id, student_id, activity_id, parent_attempt_id, started_at, status, attempt_count) VALUES (?, ?, ?, ?, ?, ?, ?)",
		att.ID, att.StudentID, att.ActivityID, att.ParentAttemptID, att.StartedAt, att.Status, att.AttemptCount,
	)
	return err
}

func (r *Repository) GetAttempt(id string) (*models.ActivityAttempt, error) {
	query := `
		SELECT id, student_id, activity_id, parent_attempt_id, started_at, completed_at, duration, status, attempt_count, created_at 
		FROM activity_attempts 
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var att models.ActivityAttempt
	var parentAttemptID sql.NullString

	if err := row.Scan(
		&att.ID,
		&att.StudentID,
		&att.ActivityID,
		&parentAttemptID,
		&att.StartedAt,
		&att.CompletedAt,
		&att.Duration,
		&att.Status,
		&att.AttemptCount,
		&att.CreatedAt,
	); err != nil {
		return nil, err
	}

	if parentAttemptID.Valid {
		att.ParentAttemptID = &parentAttemptID.String
	} else {
		att.ParentAttemptID = nil
	}

	return &att, nil
}

func (r *Repository) UpdateAttempt(id string, status string, completedAt time.Time, duration int64, attemptCount int) error {
	_, err := r.db.Exec(
		"UPDATE activity_attempts SET status = ?, completed_at = ?, duration = ?, attempt_count = ? WHERE id = ?",
		status, completedAt, duration, attemptCount, id,
	)
	return err
}

// UpdateActivity updates an activity. A blank correct_answer keeps the
// stored one, because the listing endpoint never returns answers.
func (r *Repository) UpdateActivity(a models.Activity) error {
	res, err := r.db.Exec(`
		UPDATE activities
		SET title = ?, description = ?, content = ?,
		    correct_answer = CASE WHEN ? = '' THEN correct_answer ELSE ? END,
		    topic = ?, difficulty = ?, estimated_duration_minutes = ?
		WHERE id = ?`,
		a.Title, a.Description, a.Content, a.CorrectAnswer, a.CorrectAnswer,
		a.Topic, a.Difficulty, a.EstimatedDurationMinutes, a.ID,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) RecordSignal(sig models.EngagementSignal) error {
	_, err := r.db.Exec(
		"INSERT INTO engagement_signals (id, student_id, attempt_id, signal_type, signal_value) VALUES (?, ?, ?, ?, ?)",
		sig.ID, sig.StudentID, sig.AttemptID, sig.SignalType, sig.SignalValue,
	)
	return err
}

func (r *Repository) GetStudentHistory(studentID string) ([]models.StudentHistoryItem, error) {
	query := `
		SELECT 
			a.id, a.activity_id, act.title, a.status, a.started_at, a.completed_at, a.duration, a.attempt_count,
			(SELECT COUNT(*) FROM engagement_signals s WHERE s.attempt_id = a.id) as signals_count
		FROM activity_attempts a
		JOIN activities act ON act.id = a.activity_id
		WHERE a.student_id = ?
		ORDER BY a.started_at DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.StudentHistoryItem
	for rows.Next() {
		var item models.StudentHistoryItem
		if err := rows.Scan(&item.AttemptID, &item.ActivityID, &item.ActivityTitle, &item.Status, &item.StartedAt, &item.CompletedAt, &item.Duration, &item.AttemptCount, &item.SignalsCount); err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, nil
}
