package repository

import (
	"audio-processor/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteAudioRepository struct {
	db *sql.DB
}

func NewSQLiteAudioRepository(dbPath string) (*SQLiteAudioRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &SQLiteAudioRepository{
		db: db,
	}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS audios (
		id TEXT PRIMARY KEY,
		original_name TEXT NOT NULL,
		original_format TEXT NOT NULL,
		storage_path TEXT,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		error TEXT
	)`

	_, err := db.Exec(query)
	return err
}

func (r *SQLiteAudioRepository) Store(ctx context.Context, audio *entity.Audio) error {
	query := `
	INSERT INTO audios (
		id, original_name, original_format, storage_path, 
		status, created_at, updated_at, error
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		audio.ID,
		audio.OriginalName,
		audio.OriginalFormat,
		audio.StoragePath,
		audio.Status,
		audio.CreatedAt,
		audio.UpdatedAt,
		audio.Error,
	)

	if err != nil {
		return fmt.Errorf("failed to store audio: %v", err)
	}

	return nil
}

func (r *SQLiteAudioRepository) GetByID(ctx context.Context, id string) (*entity.Audio, error) {
	query := `
	SELECT id, original_name, original_format, storage_path, 
	       status, created_at, updated_at, error
	FROM audios WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	audio := &entity.Audio{}
	var createdAt, updatedAt string

	err := row.Scan(
		&audio.ID,
		&audio.OriginalName,
		&audio.OriginalFormat,
		&audio.StoragePath,
		&audio.Status,
		&createdAt,
		&updatedAt,
		&audio.Error,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("audio not found: %s", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan audio: %v", err)
	}

	// Parse time strings
	audio.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	audio.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return audio, nil
}

func (r *SQLiteAudioRepository) Update(ctx context.Context, audio *entity.Audio) error {
	query := `
	UPDATE audios 
	SET original_name = ?,
		original_format = ?,
		storage_path = ?,
		status = ?,
		updated_at = ?,
		error = ?
	WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		audio.OriginalName,
		audio.OriginalFormat,
		audio.StoragePath,
		audio.Status,
		time.Now(),
		audio.Error,
		audio.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update audio: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("audio not found: %s", audio.ID)
	}

	return nil
}

func (r *SQLiteAudioRepository) Close() error {
	return r.db.Close()
}
