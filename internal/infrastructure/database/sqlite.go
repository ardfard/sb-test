package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes and returns a new SQLite database connection
func InitDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	return db, nil
}

// createTable creates the audios table if it doesn't exist.
func createTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS audios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_name TEXT NOT NULL,
		current_format TEXT NOT NULL,
		storage_path TEXT,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		error TEXT,
		user_id INTEGER NOT NULL,
		phrase_id INTEGER NOT NULL
	);`

	query += `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	query += `
	CREATE TABLE IF NOT EXISTS phrases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phrase TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	// Create a unique constraint on user_id and phrase_id
	query += `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_user_phrase ON audios (user_id, phrase_id);`

	_, err := db.Exec(query)
	return err
}
