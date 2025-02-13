package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteAudioRepository struct {
	db *sqlx.DB
}

func NewSQLiteAudioRepository(dbPath string) (*SQLiteAudioRepository, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &SQLiteAudioRepository{
		db: db,
	}, nil
}

func createTable(db *sqlx.DB) error {
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
	) VALUES (:id, :original_name, :original_format, :storage_path, :status, :created_at, :updated_at, :error)`
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              audio.ID,
		"original_name":   audio.OriginalName,
		"original_format": audio.OriginalFormat,
		"storage_path":    audio.StoragePath,
		"status":          audio.Status,
		"created_at":      audio.CreatedAt.Format(time.RFC3339),
		"updated_at":      audio.UpdatedAt.Format(time.RFC3339),
		"error":           audio.Error,
	})
	if err != nil {
		return fmt.Errorf("failed to store audio: %v", err)
	}

	return nil
}

func (r *SQLiteAudioRepository) GetByID(ctx context.Context, id string) (*entity.Audio, error) {
	query := `
	SELECT id, original_name, original_format, storage_path, status, created_at, updated_at, error
	FROM audios WHERE id = ?`
	audio := &entity.Audio{}
	if err := r.db.GetContext(ctx, audio, query, id); err != nil {
		return nil, fmt.Errorf("failed to get audio: %v", err)
	}
	return audio, nil
}

func (r *SQLiteAudioRepository) Update(ctx context.Context, audio *entity.Audio) error {
	query := `
	UPDATE audios 
	SET original_name = :original_name,
		original_format = :original_format,
		storage_path = :storage_path,
		status = :status,
		updated_at = :updated_at,
		error = :error
	WHERE id = :id`
	// update the updated timestamp
	audio.UpdatedAt = time.Now()
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":              audio.ID,
		"original_name":   audio.OriginalName,
		"original_format": audio.OriginalFormat,
		"storage_path":    audio.StoragePath,
		"status":          audio.Status,
		"updated_at":      audio.UpdatedAt.Format(time.RFC3339),
		"error":           audio.Error,
	})
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
