package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type migration struct {
	version int
	stmts   []string
}

var migrations = []migration{
	{
		version: 1,
		stmts: []string{
			`CREATE TABLE settings (
				key TEXT PRIMARY KEY,
				value TEXT NOT NULL
			)`,
			`CREATE TABLE providers (
				id TEXT PRIMARY KEY,
				kind TEXT NOT NULL CHECK (kind IN ('openrouter', 'openai', 'anthropic', 'deepseek', 'custom')),
				label TEXT NOT NULL,
				base_url TEXT NOT NULL DEFAULT '',
				model TEXT NOT NULL DEFAULT '',
				selected INTEGER NOT NULL DEFAULT 0 CHECK (selected IN (0, 1)),
				created_at TEXT NOT NULL,
				updated_at TEXT NOT NULL
			)`,
			`CREATE UNIQUE INDEX providers_one_selected ON providers (selected) WHERE selected = 1`,
			`CREATE TABLE recent_workbooks (
				path TEXT PRIMARY KEY,
				opened_at TEXT NOT NULL
			)`,
			`CREATE TABLE conversations (
				id INTEGER PRIMARY KEY,
				workbook_path TEXT NOT NULL DEFAULT '',
				role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
				content TEXT NOT NULL,
				created_at TEXT NOT NULL
			)`,
			`CREATE TABLE operations (
				id INTEGER PRIMARY KEY,
				workbook_path TEXT NOT NULL DEFAULT '',
				name TEXT NOT NULL,
				summary TEXT NOT NULL,
				status TEXT NOT NULL CHECK (status IN ('previewed', 'applied', 'cancelled', 'failed')),
				created_at TEXT NOT NULL
			)`,
		},
	},
}

func Migrate(ctx context.Context, db *sql.DB) error {
	if ctx == nil {
		return fmt.Errorf("migrate database: context is nil")
	}
	if db == nil {
		return fmt.Errorf("migrate database: database is nil")
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	for _, m := range migrations {
		applied, err := migrationApplied(ctx, db, m.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := applyMigration(ctx, db, m); err != nil {
			return err
		}
	}
	return nil
}

func migrationApplied(ctx context.Context, db *sql.DB, version int) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check migration %d: %w", version, err)
	}
	return count > 0, nil
}

func applyMigration(ctx context.Context, db *sql.DB, m migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %d: %w", m.version, err)
	}
	defer tx.Rollback()

	for _, stmt := range m.stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("apply migration %d: %w", m.version, err)
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, m.version, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("record migration %d: %w", m.version, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %d: %w", m.version, err)
	}
	return nil
}
