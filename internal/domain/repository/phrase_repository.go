package repository

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
)

type PhraseRepository interface {
	Create(ctx context.Context, phrase *entity.Phrase) error
	GetByID(ctx context.Context, id uint) (*entity.Phrase, error)
}
