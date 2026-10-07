package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// List of ordered schema migrations.
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		args TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`,
}

func Open(path string) (*sql.DB, error) {
	// Enable WAL mode and set a 5-second busy timeout for concurrent safety
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", path)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sqlite DSN: %w", err)
	}

	// Test the actual connection/file access
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open sqlite database at %s: %w", path, err)
	}

	// Recommended for SQLite connection pool limits
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func Migrate(database *sql.DB) error {
	for _, migration := range migrations {
		if _, err := database.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
