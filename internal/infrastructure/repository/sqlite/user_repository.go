package sqlite

import (
	"context"
	"fmt"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, name, created_at, updated_at`
	var createdUser entity.User
	err := r.db.GetContext(ctx, &createdUser, query, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &createdUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	query := `SELECT id, name, created_at, updated_at FROM users WHERE id = ?`
	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
