package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/ardfard/sb-test/pkg/projectpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDownloadAudioUseCase_Download(t *testing.T) {
	testFile := filepath.Join(projectpath.RootProject, "tests", "fixtures", "test.wav")

	tests := []struct {
		name          string
		userID        uint
		phraseID      uint
		format        string
		setupMocks    func(*repoMocks.MockAudioRepository, *storageMocks.MockStorage, *repoMocks.MockUserRepository, *repoMocks.MockPhraseRepository)
		expectedError bool
		checkResult   func(*testing.T, io.ReadCloser, error)
	}{
		{
			name:     "successful download of completed audio",
			userID:   1,
			phraseID: 1,
			format:   "m4a",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				audio := &entity.Audio{
					ID:            1,
					Status:        entity.AudioStatusCompleted,
					StoragePath:   fmt.Sprintf("%s/converted/1.wav", basePath),
					CurrentFormat: "wav",
				}
				inputReader, err := os.Open(testFile)
				if err != nil {
					t.Fatal(err)
				}

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
				repo.On("GetByUserIDAndPhraseID", mock.Anything, uint(1), uint(1)).Return(audio, nil)
				storage.On("Download", mock.Anything, fmt.Sprintf("%s/converted/1.wav", basePath)).
					Return(inputReader, nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.NoError(t, err)
				// check content of reader not nil
				content, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
			},
		},
		{
			name:     "download audio with same format",
			userID:   1,
			phraseID: 1,
			format:   "wav",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				audio := &entity.Audio{
					ID:            1,
					Status:        entity.AudioStatusPending,
					StoragePath:   fmt.Sprintf("%s/original/1.wav", basePath),
					CurrentFormat: "wav",
				}
				inputReader, err := os.Open(testFile)
				if err != nil {
					t.Fatal(err)
				}

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
				repo.On("GetByUserIDAndPhraseID", mock.Anything, uint(1), uint(1)).Return(audio, nil)
				storage.On("Download", mock.Anything, fmt.Sprintf("%s/original/1.wav", basePath)).
					Return(inputReader, nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				// check content of reader not nil
				content, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
			},
		},
		{
			name:     "download pending audio with format conversion",
			userID:   1,
			phraseID: 1,
			format:   "m4a",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				audio := &entity.Audio{
					ID:            1,
					Status:        entity.AudioStatusPending,
					StoragePath:   fmt.Sprintf("%s/original/1.wav", basePath),
					CurrentFormat: "wav",
				}

				inputReader, err := os.Open(testFile)
				if err != nil {
					t.Fatal(err)
				}

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
				repo.On("GetByUserIDAndPhraseID", mock.Anything, uint(1), uint(1)).Return(audio, nil)
				storage.On("Download", mock.Anything, fmt.Sprintf("%s/original/1.wav", basePath)).
					Return(inputReader, nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				// check content of reader not nil
				content, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
			},
		},
		{
			name:     "audio not found",
			userID:   999,
			phraseID: 1,
			format:   "wav",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				userRepo.On("GetByID", mock.Anything, uint(999)).Return(&entity.User{ID: 999}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{ID: 1}, nil)
				repo.On("GetByUserIDAndPhraseID", mock.Anything, uint(999), uint(1)).Return(nil, assert.AnError)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
			},
		},
		{
			name:     "user not found",
			userID:   999,
			phraseID: 1,
			format:   "wav",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				userRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, assert.AnError)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:     "phrase not found",
			userID:   1,
			phraseID: 999,
			format:   "wav",
			setupMocks: func(repo *repoMocks.MockAudioRepository, storage *storageMocks.MockStorage, userRepo *repoMocks.MockUserRepository, phraseRepo *repoMocks.MockPhraseRepository) {
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{ID: 1}, nil)
				phraseRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, assert.AnError)
			},
			expectedError: true,
			checkResult: func(t *testing.T, reader io.ReadCloser, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoMocks.MockAudioRepository{}
			storage := &storageMocks.MockStorage{}
			audioConverter := converter.NewAudioConverter()
			userRepo := &repoMocks.MockUserRepository{}
			phraseRepo := &repoMocks.MockPhraseRepository{}

			tt.setupMocks(repo, storage, userRepo, phraseRepo)

			uc := NewDownloadAudioUseCase(repo, storage, audioConverter, userRepo, phraseRepo)
			reader, err := uc.Download(context.Background(), tt.userID, tt.phraseID, tt.format)

			tt.checkResult(t, reader, err)

			repo.AssertExpectations(t)
			storage.AssertExpectations(t)

		})
	}
}
