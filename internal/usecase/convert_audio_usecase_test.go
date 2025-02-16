package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/ardfard/sb-test/pkg/projectpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConvertAudioUseCase_Convert(t *testing.T) {
	testFile := filepath.Join(projectpath.RootProject, "tests", "fixtures", "test.m4a")
	reader, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		audioID       uint
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage)
		expectedError bool
	}{
		{
			name:    "successful conversion",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				audio := &entity.Audio{
					ID:            1,
					CurrentFormat: "m4a",
					Status:        entity.AudioStatusPending,
					StoragePath:   fmt.Sprintf("%s/original/%d.m4a", basePath, 1),
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.Status == entity.AudioStatusConverting
				})).Return(nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.Status == entity.AudioStatusCompleted
				})).Return(nil)

				storage.On("Download", mock.Anything, fmt.Sprintf("%s/original/%d.m4a", basePath, 1)).
					Return(reader, nil)
				storage.On("Upload", mock.Anything, mock.MatchedBy(func(path string) bool {
					return strings.HasSuffix(path, ".wav")
				}), mock.Anything).Return(nil)
				storage.On("Delete", mock.Anything, fmt.Sprintf("%s/original/%d.m4a", basePath, 1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "repository error",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(nil, assert.AnError)
			},
			expectedError: true,
		},
		{
			name:    "storage download error",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage) {
				audio := &entity.Audio{
					ID:            1,
					CurrentFormat: "mp3",
					Status:        entity.AudioStatusPending,
					StoragePath:   "original/1.mp3",
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
				storage.On("Download", mock.Anything, "original/1.mp3").Return(nil, assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}
			audioConverter := converter.NewAudioConverter()

			tt.setupMocks(repo, storage)

			uc := NewConvertAudioUseCase(repo, storage, audioConverter)
			err := uc.Convert(context.Background(), tt.audioID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			storage.AssertExpectations(t)
		})
	}
}
