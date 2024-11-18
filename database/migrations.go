package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

// List of migrations that should be applied.
// Migration ID = index + 1.
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS admins (
        id INTEGER PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        bcrypt TEXT NOT NULL
    )`,
	`CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY,
        name TEXT
    )`,
	`CREATE TABLE IF NOT EXISTS links (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        href TEXT NOT NULL,
        group_id INTEGER NOT NULL,
        icon TEXT,
        CONSTRAINT fk_groups_links
            FOREIGN KEY (group_id)
            REFERENCES groups(id)
            ON DELETE CASCADE
    )`,
}

func ApplyMigrations() error {
	// Create a table to record applied migrations and retrieve the saved data.
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (
        version INTEGER NOT NULL DEFAULT 0
    )`)
	if err != nil {
		return err
	}

	var currentVersion int
	err = DB.
		QueryRow("SELECT version FROM migrations").
		Scan(&currentVersion)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// The table is empty, create a record.
		_, err = DB.Exec("INSERT INTO migrations (version) VALUES (0)")
		if err != nil {
			return err
		}
	}

	// Apply all migrations.
	for i, migration := range migrations {
		migrationID := i + 1
		if migrationID <= currentVersion {
			continue
		}

		if err := applyMigration(migrationID, migration); err != nil {
			return fmt.Errorf("migration #%d: %w", migrationID, err)
		}

		slog.Info("migration has been applied", "id", migrationID)
	}

	return nil
}

func applyMigration(migrationID int, query string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("error when applying migration: %w", err)
	}

	if _, err := tx.Exec("UPDATE migrations SET version = ?", migrationID); err != nil {
		return fmt.Errorf("error when updating schema version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
