package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestMigrateCreatesTablesOnce(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("second Migrate() error = %v", err)
	}

	var applied int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM schema_migrations`).Scan(&applied); err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if applied != len(migrations) {
		t.Fatalf("applied migrations = %d, want %d", applied, len(migrations))
	}

	for _, table := range []string{"settings", "providers", "recent_workbooks", "conversations", "operations"} {
		var name string
		err := db.QueryRowContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&name)
		if err != nil {
			t.Fatalf("table %s: %v", table, err)
		}
	}
}

func TestOnlyOneProviderCanBeSelected(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	insert := `INSERT INTO providers (id, kind, label, base_url, model, selected, created_at, updated_at)
		VALUES (?, 'openrouter', 'OpenRouter', '', 'model', ?, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`
	if _, err := db.ExecContext(ctx, insert, "a", 1); err != nil {
		t.Fatalf("insert selected provider: %v", err)
	}
	if _, err := db.ExecContext(ctx, insert, "b", 1); err == nil {
		t.Fatal("second selected provider was inserted")
	}
}

func TestMigrateRejectsNilDatabase(t *testing.T) {
	if err := Migrate(context.Background(), nil); err == nil {
		t.Fatal("Migrate() error = nil, want an error")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Open(context.Background(), filepath.Join(t.TempDir(), "cellpilot.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
	return db
}
