package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	queueMocks "github.com/ardfard/sb-test/internal/infrastructure/queue/mocks"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadAudioUseCase_Upload(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage, *queueMocks.MockTaskQueue, *repoMocks.MockUserRepository, *repoMocks.MockPhraseRepository)
		expectedError bool
	}{
		{
			name:     "successful upload",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				repo.On("Store", mock.Anything, mock.MatchedBy(func(audio *entity.Audio) bool {
					return audio.OriginalName == "test.mp3" && audio.CurrentFormat == ".mp3"
				})).Return(&entity.Audio{ID: 1, OriginalName: "test.mp3", CurrentFormat: ".mp3"}, nil)

				storage.On("Upload", mock.Anything, mock.MatchedBy(func(path string) bool {
					return strings.HasPrefix(path, fmt.Sprintf("%s/original/", basePath)) && strings.HasSuffix(path, ".mp3")
				}), mock.Anything).Return(nil)

				queue.On("Enqueue", mock.Anything, mock.AnythingOfType("uint")).Return(nil)

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
			},
			expectedError: false,
		},
		{
			name:     "repository error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				repo.On("Store", mock.Anything, mock.Anything).Return(nil, assert.AnError)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
			},
			expectedError: true,
		},
		{
			name:     "storage error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				repo.On("Store", mock.Anything, mock.Anything).Return(&entity.Audio{ID: 1, OriginalName: "test.mp3", CurrentFormat: ".mp3"}, nil)
				storage.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
			},
			expectedError: true,
		},
		{
			name:     "queue error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				repo.On("Store", mock.Anything, mock.Anything).Return(&entity.Audio{ID: 1, OriginalName: "test.mp3", CurrentFormat: ".mp3"}, nil)
				storage.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				queue.On("Enqueue", mock.Anything, mock.Anything).Return(assert.AnError)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
			},
			expectedError: true,
		},
		{
			name:     "user not found",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(nil, assert.AnError)
			},
			expectedError: true,
		},
		{
			name:     "phrase not found",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(nil, assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}
			queue := &queueMocks.MockTaskQueue{}
			userRepo := &repoMocks.MockUserRepository{}
			phraseRepo := &repoMocks.MockPhraseRepository{}

			tt.setupMocks(repo, storage, queue, userRepo, phraseRepo)

			uc := NewUploadAudioUseCase(repo, storage, queue, userRepo, phraseRepo)
			content := strings.NewReader("test content")
			audio, err := uc.Upload(context.Background(), tt.filename, content, 1, 1)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, audio)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, audio)
				assert.Equal(t, tt.filename, audio.OriginalName)
			}

			repo.AssertExpectations(t)
			storage.AssertExpectations(t)
			queue.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			phraseRepo.AssertExpectations(t)
		})
	}
}
