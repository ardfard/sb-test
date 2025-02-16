package usecase

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
)

type CreateUserUsecase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUsecase(userRepository repository.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{userRepository: userRepository}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, user *entity.User) error {
	return u.userRepository.Create(ctx, user)
}
