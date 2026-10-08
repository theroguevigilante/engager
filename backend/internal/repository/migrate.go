package repository

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Migrate applies every *.sql file in dir (lexical order) to the database.
// It is safe to run on every start and on databases created by earlier
// versions: CREATE statements use IF NOT EXISTS, and ALTER TABLE ADD COLUMN
// statements whose column already exists are skipped.
func (r *Repository) Migrate(fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		raw, err := fs.ReadFile(fsys, path.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		for _, stmt := range splitStatements(string(raw)) {
			if _, err := r.db.Exec(stmt); err != nil {
				if strings.Contains(err.Error(), "duplicate column name") {
					continue
				}
				return fmt.Errorf("migration %s: %w", name, err)
			}
		}
	}
	return nil
}

// splitStatements strips "--" line comments and splits on semicolons.
func splitStatements(sqlText string) []string {
	var cleaned []string
	for _, line := range strings.Split(sqlText, "\n") {
		if i := strings.Index(line, "--"); i >= 0 {
			line = line[:i]
		}
		cleaned = append(cleaned, line)
	}
	var out []string
	for _, s := range strings.Split(strings.Join(cleaned, "\n"), ";") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
