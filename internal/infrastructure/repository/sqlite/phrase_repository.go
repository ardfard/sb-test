package sqlite

import (
	"context"
	"fmt"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type PhraseRepository struct {
	db *sqlx.DB
}

func NewPhraseRepository(db *sqlx.DB) (*PhraseRepository, error) {
	return &PhraseRepository{db: db}, nil
}

func (r *PhraseRepository) Create(ctx context.Context, phrase *entity.Phrase) error {
	query := `INSERT INTO phrases (phrase, created_at, updated_at) VALUES (:phrase, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"phrase":     phrase.Phrase,
		"created_at": phrase.CreatedAt,
		"updated_at": phrase.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create phrase: %w", err)
	}
	return nil
}

func (r *PhraseRepository) GetByID(ctx context.Context, id uint) (*entity.Phrase, error) {
	query := `SELECT id, phrase, created_at, updated_at FROM phrases WHERE id = ?`
	var phrase entity.Phrase
	err := r.db.GetContext(ctx, &phrase, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get phrase: %w", err)
	}
	return &phrase, nil
}
