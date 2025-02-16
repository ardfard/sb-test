package worker

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/queue"
	convertermocks "github.com/ardfard/sb-test/internal/infrastructure/converter/mocks"
	queuemocks "github.com/ardfard/sb-test/internal/infrastructure/queue/mocks"
	repomocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storagemocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReadCloser struct {
	io.Reader
	closeFunc func() error
}

func (m mockReadCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func newMockReadCloser(s string) io.ReadCloser {
	return mockReadCloser{Reader: strings.NewReader(s)}
}

func TestConversionWorker_ProcessNextMessage(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*queuemocks.MockTaskQueue, *repomocks.MockAudioRepository, *storagemocks.MockStorage, *convertermocks.MockAudioConverter)
		expectedError  error
		expectedErrMsg string
	}{
		{
			name: "Success",
			setupMocks: func(mockQueue *queuemocks.MockTaskQueue, mockRepo *repomocks.MockAudioRepository, mockStorage *storagemocks.MockStorage, mockConverter *convertermocks.MockAudioConverter) {
				task := &queue.Task{
					ID:      "1",
					Payload: uint(123),
				}

				audio := &entity.Audio{
					ID:            123,
					UserID:        1,
					PhraseID:      1,
					CurrentFormat: "mp3",
					Status:        "pending",
					StoragePath:   "test/path",
				}

				// Mock repository calls
				mockRepo.On("GetByID", mock.Anything, uint(123)).Return(audio, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.ID == 123 && a.Status == "converting"
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.ID == 123 && a.Status == "completed"
				})).Return(nil).Once()

				// Mock storage calls
				mockStorage.On("Download", mock.Anything, audio.StoragePath).Return(newMockReadCloser("test data"), nil)
				mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockStorage.On("Delete", mock.Anything, audio.StoragePath).Return(nil)

				// Mock converter calls
				mockConverter.On("ConvertFromReader", mock.Anything, mock.Anything, audio.CurrentFormat, "wav").Return(newMockReadCloser("converted data"), nil)

				// Mock queue calls
				mockQueue.On("Dequeue", mock.Anything).Return(task, nil)
				mockQueue.On("Complete", mock.Anything, task.ID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Dequeue Error",
			setupMocks: func(mockQueue *queuemocks.MockTaskQueue, mockRepo *repomocks.MockAudioRepository, mockStorage *storagemocks.MockStorage, mockConverter *convertermocks.MockAudioConverter) {
				expectedErr := errors.New("dequeue error")
				mockQueue.On("Dequeue", mock.Anything).Return(nil, expectedErr)
			},
			expectedError:  errors.New("failed to dequeue message: dequeue error"),
			expectedErrMsg: "failed to dequeue message",
		},
		{
			name: "Conversion Error",
			setupMocks: func(mockQueue *queuemocks.MockTaskQueue, mockRepo *repomocks.MockAudioRepository, mockStorage *storagemocks.MockStorage, mockConverter *convertermocks.MockAudioConverter) {
				task := &queue.Task{
					ID:      "1",
					Payload: uint(123),
				}

				audio := &entity.Audio{
					ID:            123,
					UserID:        1,
					PhraseID:      1,
					CurrentFormat: "mp3",
					Status:        "pending",
					StoragePath:   "test/path",
				}

				conversionErr := errors.New("conversion error")

				// Mock repository calls
				mockRepo.On("GetByID", mock.Anything, uint(123)).Return(audio, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.ID == 123 && a.Status == "converting"
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *entity.Audio) bool {
					return a.ID == 123 && a.Status == "failed"
				})).Return(nil).Once()

				// Mock storage calls
				mockStorage.On("Download", mock.Anything, audio.StoragePath).Return(newMockReadCloser("test data"), nil)

				// Mock converter calls
				mockConverter.On("ConvertFromReader", mock.Anything, mock.Anything, audio.CurrentFormat, "wav").Return(nil, conversionErr)

				// Mock queue calls
				mockQueue.On("Dequeue", mock.Anything).Return(task, nil)
				mockQueue.On("Fail", mock.Anything, task.ID, mock.Anything).Return(nil)
			},
			expectedError:  errors.New("failed to convert audio: conversion error"),
			expectedErrMsg: "failed to convert audio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockQueue := queuemocks.NewMockTaskQueue(t)
			mockRepo := repomocks.NewMockAudioRepository(t)
			mockStorage := storagemocks.NewMockStorage(t)
			mockConverter := convertermocks.NewMockAudioConverter(t)
			useCase := usecase.NewConvertAudioUseCase(mockRepo, mockStorage, mockConverter)

			// Setup test case mocks
			tt.setupMocks(mockQueue, mockRepo, mockStorage, mockConverter)

			// Create worker
			worker := NewConversionWorker(mockQueue, useCase)

			// Run test
			err := worker.processNextMessage()

			// Cleanup
			os.RemoveAll("test")

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}

			mockQueue.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
			mockConverter.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestConversionWorker_StartStop(t *testing.T) {
	// Setup mocks
	mockQueue := queuemocks.NewMockTaskQueue(t)
	mockAudioRepository := repomocks.NewMockAudioRepository(t)
	mockStorage := storagemocks.NewMockStorage(t)
	mockConverter := convertermocks.NewMockAudioConverter(t)
	useCase := usecase.NewConvertAudioUseCase(mockAudioRepository, mockStorage, mockConverter)

	// Create worker
	worker := NewConversionWorker(mockQueue, useCase)

	// Start worker
	worker.Start()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop worker
	worker.Stop()

	// Verify it stops (this is a bit tricky to test thoroughly)
	// In a real scenario, you might want to add more sophisticated verification
	assert.NotNil(t, worker.stopChan)
}
