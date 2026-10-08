package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (r *Repository) SaveCheckin(studentID, need string) error {
	_, err := r.db.Exec("INSERT INTO student_checkins (id, student_id, need) VALUES (?, ?, ?)",
		uuid.New().String(), studentID, need)
	return err
}

// LatestCheckin returns the student's most recent need from the last 24 hours ("" if none).
func (r *Repository) LatestCheckin(studentID string) (string, error) {
	var need string
	err := r.db.QueryRow(`SELECT need FROM student_checkins
		WHERE student_id = ? AND created_at >= datetime('now', '-24 hours')
		ORDER BY created_at DESC, rowid DESC LIMIT 1`, studentID).Scan(&need)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return need, err
}

// PulseCounts returns how many distinct students reported each need in the last 24 hours.
func (r *Repository) PulseCounts() (map[string]int, error) {
	rows, err := r.db.Query(`SELECT need, COUNT(DISTINCT student_id) FROM student_checkins
		WHERE created_at >= datetime('now', '-24 hours') GROUP BY need`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var need string
		var n int
		if err := rows.Scan(&need, &n); err != nil {
			return nil, err
		}
		out[need] = n
	}
	return out, rows.Err()
}
