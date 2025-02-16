package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

// InitDB initializes and returns a new SQLite database connection
func InitDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := createTable(db); err != nil {
		db.Close() // Close the db connection if table creation fails
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	return db, nil
}

// createTable creates the audios table if it doesn't exist.
func createTable(db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	schemaSQL, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}
