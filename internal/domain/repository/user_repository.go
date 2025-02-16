package repository

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
}
