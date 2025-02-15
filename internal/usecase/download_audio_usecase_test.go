package usecase

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDownloadAudioUseCase_Download(t *testing.T) {
	tests := []struct {
		name          string
		audioID       uint
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage)
		expectedError bool
		checkResult   func(*testing.T, io.ReadCloser, error)
	}{
		{
			name:    "successful download",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				audio := &entity.Audio{
					ID:          1,
					Status:      entity.AudioStatusCompleted,
					StoragePath: "converted/1.wav",
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				storage.On("Download", mock.Anything, "converted/1.wav").
					Return(mockReadCloser{strings.NewReader("test content")}, nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)

				content, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.Equal(t, "test content", string(content))
			},
		},
		{
			name:    "audio not found",
			audioID: 999,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				repo.On("GetByID", mock.Anything, uint(999)).Return(nil, assert.AnError)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
			},
		},
		{
			name:    "audio not ready",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				audio := &entity.Audio{
					ID:          1,
					Status:      entity.AudioStatusConverting,
					StoragePath: "converted/1.wav",
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
				assert.Contains(t, err.Error(), "not ready for download")
			},
		},
		{
			name:    "storage error",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				audio := &entity.Audio{
					ID:          1,
					Status:      entity.AudioStatusCompleted,
					StoragePath: "converted/1.wav",
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				storage.On("Download", mock.Anything, "converted/1.wav").
					Return(nil, assert.AnError)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
				assert.Contains(t, err.Error(), "failed to download file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}

			tt.setupMocks(repo, storage)

			uc := NewDownloadAudioUseCase(repo, storage)
			reader, err := uc.Download(context.Background(), tt.audioID)

			tt.checkResult(t, reader, err)

			repo.AssertExpectations(t)
			storage.AssertExpectations(t)

			if reader != nil {
				reader.Close()
			}
		})
	}
}
