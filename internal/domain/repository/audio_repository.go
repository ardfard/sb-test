package repository

import (
	"audio-processor/internal/domain/entity"
	"context"
)

type AudioRepository interface {
	Store(ctx context.Context, audio *entity.Audio) error
	GetByID(ctx context.Context, id string) (*entity.Audio, error)
	Update(ctx context.Context, audio *entity.Audio) error
}
