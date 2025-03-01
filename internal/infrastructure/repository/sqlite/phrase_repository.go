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

func (r *PhraseRepository) Create(ctx context.Context, phrase *entity.Phrase) (*entity.Phrase, error) {
	query := `INSERT INTO phrases (phrase, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, phrase, created_at, updated_at`
	var createdPhrase entity.Phrase
	err := r.db.GetContext(ctx, &createdPhrase, query, phrase.Phrase, phrase.CreatedAt, phrase.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create phrase: %w", err)
	}
	return &createdPhrase, nil
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
