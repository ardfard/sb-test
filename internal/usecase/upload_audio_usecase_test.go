package usecase

import (
	"context"
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
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage, *queueMocks.MockTaskQueue)
		expectedError bool
	}{
		{
			name:     "successful upload",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue) {
				repo.On("Store", mock.Anything, mock.MatchedBy(func(audio *entity.Audio) bool {
					return audio.OriginalName == "test.mp3" && audio.OriginalFormat == ".mp3"
				})).Return(nil)

				storage.On("Upload", mock.Anything, mock.MatchedBy(func(path string) bool {
					return strings.HasPrefix(path, "original/") && strings.HasSuffix(path, ".mp3")
				}), mock.Anything).Return(nil)

				queue.On("Enqueue", mock.Anything, mock.AnythingOfType("uint")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "repository error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue) {
				repo.On("Store", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedError: true,
		},
		{
			name:     "storage error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue) {
				repo.On("Store", mock.Anything, mock.Anything).Return(nil)
				storage.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedError: true,
		},
		{
			name:     "queue error",
			filename: "test.mp3",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, queue *queueMocks.MockTaskQueue) {
				repo.On("Store", mock.Anything, mock.Anything).Return(nil)
				storage.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				queue.On("Enqueue", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}
			queue := &queueMocks.MockTaskQueue{}

			tt.setupMocks(repo, storage, queue)

			uc := NewUploadAudioUseCase(repo, storage, queue)
			content := strings.NewReader("test content")
			audio, err := uc.Upload(context.Background(), tt.filename, content)

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
		})
	}
}
