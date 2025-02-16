package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteAudioRepository is a repository for audio operations using SQLite.
type AudioRepository struct {
	db *sqlx.DB
}

// NewAudioRepository creates a new AudioRepository.
func NewAudioRepository(db *sqlx.DB) (*AudioRepository, error) {

	return &AudioRepository{
		db: db,
	}, nil
}

// Store stores an audio entity in the database.
func (r *AudioRepository) Store(ctx context.Context, audio *entity.Audio) error {
	query := `
	INSERT INTO audios (
		id, original_name, current_format, storage_path, 
		status, created_at, updated_at, error,
		user_id, phrase_id
	) VALUES (:id, :original_name, :current_format, :storage_path, 
		:status, :created_at, :updated_at, :error,
		:user_id, :phrase_id)`
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":             audio.ID,
		"original_name":  audio.OriginalName,
		"current_format": audio.CurrentFormat,
		"storage_path":   audio.StoragePath,
		"status":         audio.Status,
		"created_at":     audio.CreatedAt.Format(time.RFC3339),
		"updated_at":     audio.UpdatedAt.Format(time.RFC3339),
		"error":          audio.Error,
		"user_id":        audio.UserID,
		"phrase_id":      audio.PhraseID,
	})
	if err != nil {
		return fmt.Errorf("failed to store audio: %v", err)
	}
	return nil
}

// GetByID retrieves an audio entity from the database by its ID.
func (r *AudioRepository) GetByID(ctx context.Context, id uint) (*entity.Audio, error) {
	query := `
	SELECT id, original_name, current_format, storage_path, status, 
		created_at, updated_at, error, user_id, phrase_id
	FROM audios WHERE id = ?`
	audio := &entity.Audio{}
	if err := r.db.GetContext(ctx, audio, query, id); err != nil {
		return nil, fmt.Errorf("failed to get audio: %v", err)
	}
	return audio, nil
}

// Update updates an audio entity in the database.
func (r *AudioRepository) Update(ctx context.Context, audio *entity.Audio) error {
	query := `
	UPDATE audios 
	SET original_name = :original_name,
		current_format = :current_format,
		storage_path = :storage_path,
		status = :status,
		updated_at = :updated_at,
		error = :error
	WHERE id = :id`
	// update the updated timestamp
	audio.UpdatedAt = time.Now()
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":             audio.ID,
		"original_name":  audio.OriginalName,
		"current_format": audio.CurrentFormat,
		"storage_path":   audio.StoragePath,
		"status":         audio.Status,
		"updated_at":     audio.UpdatedAt.Format(time.RFC3339),
		"error":          audio.Error,
	})
	if err != nil {
		return fmt.Errorf("failed to update audio: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("audio not found: %d", audio.ID)
	}

	return nil
}

// GetByUserIDAndPhraseID retrieves an audio entity from the database by user ID and phrase ID.
func (r *AudioRepository) GetByUserIDAndPhraseID(ctx context.Context, userID uint, phraseID uint) (*entity.Audio, error) {
	query := `
	SELECT id, original_name, current_format, storage_path, status, 
		created_at, updated_at, error, user_id, phrase_id
	FROM audios WHERE user_id = ? AND phrase_id = ?`
	audio := &entity.Audio{}
	if err := r.db.GetContext(ctx, audio, query, userID, phraseID); err != nil {
		return nil, fmt.Errorf("failed to get audio: %v", err)
	}
	return audio, nil
}

// Close closes the SQLite database connection.
func (r *AudioRepository) Close() error {
	return r.db.Close()
}
