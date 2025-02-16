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

func TestCreatePhraseUseCase_Create(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		userID         uint
		mockSetup      func(*repoMocks.MockPhraseRepository)
		expectedError  error
		expectedPhrase *entity.Phrase
	}{
		{
			name:   "successful creation",
			text:   "Hello, World!",
			userID: 1,
			mockSetup: func(repo *repoMocks.MockPhraseRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(phrase *entity.Phrase) bool {
					return phrase.Phrase == "Hello, World!" && phrase.UserID == uint(1)
				})).Return(&entity.Phrase{
					ID:     1,
					UserID: 1,
					Phrase: "Hello, World!",
				}, nil)
			},
			expectedError: nil,
			expectedPhrase: &entity.Phrase{
				ID:     1,
				UserID: 1,
				Phrase: "Hello, World!",
			},
		},
		{
			name:   "repository error",
			text:   "Test Phrase",
			userID: 1,
			mockSetup: func(repo *repoMocks.MockPhraseRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(phrase *entity.Phrase) bool {
					return phrase.Phrase == "Test Phrase" && phrase.UserID == uint(1)
				})).Return(nil, errors.New("database error"))
			},
			expectedError:  errors.New("database error"),
			expectedPhrase: nil,
		},
		{
			name:   "empty text",
			text:   "",
			userID: 1,
			mockSetup: func(repo *repoMocks.MockPhraseRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(phrase *entity.Phrase) bool {
					return phrase.Phrase == "" && phrase.UserID == uint(1)
				})).Return(&entity.Phrase{
					ID:     1,
					UserID: 1,
					Phrase: "",
				}, nil)
			},
			expectedError: nil,
			expectedPhrase: &entity.Phrase{
				ID:     1,
				UserID: 1,
				Phrase: "",
			},
		},
		{
			name:   "zero user ID",
			text:   "Test Phrase",
			userID: 0,
			mockSetup: func(repo *repoMocks.MockPhraseRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(phrase *entity.Phrase) bool {
					return phrase.Phrase == "Test Phrase" && phrase.UserID == uint(0)
				})).Return(&entity.Phrase{
					ID:     1,
					UserID: 0,
					Phrase: "Test Phrase",
				}, nil)
			},
			expectedError: nil,
			expectedPhrase: &entity.Phrase{
				ID:     1,
				UserID: 0,
				Phrase: "Test Phrase",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock repository
			mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)

			// Set up mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mockPhraseRepo)
			}

			// Create use case
			useCase := usecase.NewCreatePhraseUseCase(mockPhraseRepo)

			// Execute use case
			phrase, err := useCase.Create(context.Background(), tt.text, tt.userID)

			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Check result
			assert.Equal(t, tt.expectedPhrase, phrase)

			// Verify all expectations were met
			mockPhraseRepo.AssertExpectations(t)
		})
	}
}

func TestNewCreatePhraseUseCase(t *testing.T) {
	mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)
	useCase := usecase.NewCreatePhraseUseCase(mockPhraseRepo)
	assert.NotNil(t, useCase)
}
