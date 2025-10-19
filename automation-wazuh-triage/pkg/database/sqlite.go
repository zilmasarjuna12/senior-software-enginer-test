package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitSQLite(dbPath string) (*sql.DB, error) {
	// Create database directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create closed_events table
	if err := createClosedEventsTable(db); err != nil {
		return nil, fmt.Errorf("failed to create closed_events table: %w", err)
	}

	return db, nil
}

func createClosedEventsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS closed_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id TEXT NOT NULL,
			rule_id TEXT,
			raw_event TEXT,
			reason TEXT NOT NULL,
			status TEXT NOT NULL,
			close_at DATETIME NOT NULL,
			UNIQUE(event_id)
		);
		CREATE INDEX IF NOT EXISTS idx_event_id ON closed_events(event_id);
		CREATE INDEX IF NOT EXISTS idx_close_at ON closed_events(close_at);
	`

	_, err := db.Exec(query)
	return err
}
