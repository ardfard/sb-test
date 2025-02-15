package usecase

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	converterMocks "github.com/ardfard/sb-test/internal/infrastructure/converter/mocks"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/ardfard/sb-test/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReadCloser struct {
	io.Reader
}

func (m mockReadCloser) Close() error {
	return nil
}

func TestConvertAudioUseCase_Convert(t *testing.T) {
	tests := []struct {
		name          string
		audioID       uint
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage, *converterMocks.MockAudioConverter)
		expectedError bool
	}{
		{
			name:    "successful conversion",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, converter *converterMocks.MockAudioConverter) {
				audio := &entity.Audio{
					ID:             1,
					OriginalFormat: ".mp3",
					Status:         entity.AudioStatusPending,
					StoragePath:    "original/1.mp3",
				}

				// create temporary files
				inputPath, outputPath, err := util.CreateTemporaryFiles(audio, "flac")

				os.WriteFile(inputPath, []byte("test content"), 0644)
				os.WriteFile(outputPath, []byte("test content"), 0644)

				if err != nil {
					t.Fatal(err)
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.Status == entity.AudioStatusConverting
				})).Return(nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.Status == entity.AudioStatusCompleted
				})).Return(nil)

				storage.On("Download", mock.Anything, "original/1.mp3").
					Return(mockReadCloser{strings.NewReader("test content")}, nil)
				storage.On("Upload", mock.Anything, "converted/1.wav", mock.Anything).Return(nil)

				converter.On("Convert", mock.Anything, mock.Anything, mock.Anything, "flac").Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "repository error",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, converter *converterMocks.MockAudioConverter) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(nil, assert.AnError)
			},
			expectedError: true,
		},
		{
			name:    "conversion error",
			audioID: 1,
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, converter *converterMocks.MockAudioConverter) {
				audio := &entity.Audio{
					ID:             1,
					OriginalFormat: ".mp3",
					Status:         entity.AudioStatusPending,
					StoragePath:    "original/1.mp3",
				}

				repo.On("GetByID", mock.Anything, uint(1)).Return(audio, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)

				storage.On("Download", mock.Anything, "original/1.mp3").
					Return(mockReadCloser{strings.NewReader("test content")}, nil)

				converter.On("Convert", mock.Anything, mock.Anything, mock.Anything, "flac").Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "audio_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}
			converter := &converterMocks.MockAudioConverter{}

			tt.setupMocks(repo, storage, converter)

			uc := NewConvertAudioUseCase(repo, storage, converter)
			err = uc.Convert(context.Background(), tt.audioID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			storage.AssertExpectations(t)
			converter.AssertExpectations(t)
		})
	}
}
