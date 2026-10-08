package repository

import (
	"io/fs"
	"path"
	"testing"

	"github.com/adaptive-learning/backend/migrations"
)

func TestEmbeddedMigrationsReadable(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil || len(entries) < 6 {
		t.Fatalf("entries=%d err=%v", len(entries), err)
	}
	for _, e := range entries {
		if _, err := fs.ReadFile(migrations.FS, path.Join(".", e.Name())); err != nil {
			t.Fatalf("%s: %v", e.Name(), err)
		}
	}
}
