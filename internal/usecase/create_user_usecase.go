package usecase

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUseCase(userRepository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *CreateUserUseCase) Create(ctx context.Context, name string) (*entity.User, error) {
	user := &entity.User{
		Name: name,
	}

	user, err := uc.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
