package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenCreatesDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "cellpilot.db")
	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	var one int
	if err := db.QueryRowContext(context.Background(), "SELECT 1").Scan(&one); err != nil {
		t.Fatalf("SELECT 1 error = %v", err)
	}
	if one != 1 {
		t.Fatalf("SELECT 1 = %d, want 1", one)
	}

	var foreignKeys int
	if err := db.QueryRowContext(context.Background(), "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("PRAGMA foreign_keys error = %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.IsDir() {
		t.Fatal("database path is a directory")
	}
}

func TestOpenRejectsEmptyPath(t *testing.T) {
	_, err := Open(context.Background(), " ")
	if err == nil {
		t.Fatal("Open() error = nil, want an error")
	}
}

func TestOpenRequiresContext(t *testing.T) {
	_, err := Open(nil, filepath.Join(t.TempDir(), "cellpilot.db"))
	if err == nil {
		t.Fatal("Open() error = nil, want an error")
	}
}

func TestDefaultPathUsesConfigDirectory(t *testing.T) {
	path, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath() error = %v", err)
	}
	if !strings.HasSuffix(path, filepath.Join("CellPilot", "cellpilot.db")) {
		t.Fatalf("DefaultPath() = %q, want it to end with CellPilot%ccellpilot.db", path, os.PathSeparator)
	}
}
