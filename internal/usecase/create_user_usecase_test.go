package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserUseCase_Create(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		mockSetup     func(*repoMocks.MockUserRepository)
		expectedError error
		expectedUser  *entity.User
	}{
		{
			name:     "successful creation",
			userName: "John Doe",
			mockSetup: func(repo *repoMocks.MockUserRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == "John Doe"
				})).Return(&entity.User{
					ID:   1,
					Name: "John Doe",
				}, nil)
			},
			expectedError: nil,
			expectedUser: &entity.User{
				ID:   1,
				Name: "John Doe",
			},
		},
		{
			name:     "repository error",
			userName: "Jane Doe",
			mockSetup: func(repo *repoMocks.MockUserRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == "Jane Doe"
				})).Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectedUser:  nil,
		},
		{
			name:     "empty name",
			userName: "",
			mockSetup: func(repo *repoMocks.MockUserRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == ""
				})).Return(&entity.User{
					ID:   1,
					Name: "",
				}, nil)
			},
			expectedError: nil,
			expectedUser: &entity.User{
				ID:   1,
				Name: "",
			},
		},
		{
			name:     "very long name",
			userName: "This is a very long name that might exceed database limits",
			mockSetup: func(repo *repoMocks.MockUserRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == "This is a very long name that might exceed database limits"
				})).Return(&entity.User{
					ID:   1,
					Name: "This is a very long name that might exceed database limits",
				}, nil)
			},
			expectedError: nil,
			expectedUser: &entity.User{
				ID:   1,
				Name: "This is a very long name that might exceed database limits",
			},
		},
		{
			name:     "special characters in name",
			userName: "User@#$%^&*()",
			mockSetup: func(repo *repoMocks.MockUserRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == "User@#$%^&*()"
				})).Return(&entity.User{
					ID:   1,
					Name: "User@#$%^&*()",
				}, nil)
			},
			expectedError: nil,
			expectedUser: &entity.User{
				ID:   1,
				Name: "User@#$%^&*()",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockUserRepo := repoMocks.NewMockUserRepository(t)

			// Set up mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mockUserRepo)
			}

			// Create use case
			useCase := usecase.NewCreateUserUseCase(mockUserRepo)

			// Execute use case
			user, err := useCase.Create(context.Background(), tt.userName)

			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Check result
			assert.Equal(t, tt.expectedUser, user)

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestNewCreateUserUseCase(t *testing.T) {
	mockUserRepo := repoMocks.NewMockUserRepository(t)
	useCase := usecase.NewCreateUserUseCase(mockUserRepo)
	assert.NotNil(t, useCase)
}
