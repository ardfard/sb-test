package repository

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
)

type AudioRepository interface {
	Store(ctx context.Context, audio *entity.Audio) (*entity.Audio, error)
	GetByID(ctx context.Context, id uint) (*entity.Audio, error)
	GetByUserIDAndPhraseID(ctx context.Context, userID uint, phraseID uint) (*entity.Audio, error)
	Update(ctx context.Context, audio *entity.Audio) error
}
