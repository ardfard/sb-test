package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/ardfard/sb-test/internal/domain/entity"
	converterMocks "github.com/ardfard/sb-test/internal/infrastructure/converter/mocks"
	queueMocks "github.com/ardfard/sb-test/internal/infrastructure/queue/mocks"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	storageMocks "github.com/ardfard/sb-test/internal/infrastructure/storage/mocks"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestAudioHandler_Upload(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		phraseID       string
		fileContent    string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful upload",
			userID:         "1",
			phraseID:       "1",
			fileContent:    "test audio content",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":     float64(1),
				"status": "pending",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			mockAudioRepo := repoMocks.NewMockAudioRepository(t)
			mockUserRepo := repoMocks.NewMockUserRepository(t)
			mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)
			mockStorage := storageMocks.NewMockStorage(t)
			mockConverter := converterMocks.NewMockAudioConverter(t)
			mockQueue := queueMocks.NewMockTaskQueue(t)

			// Set up mock expectations
			mockUserRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{
				ID:   1,
				Name: "Test User",
			}, nil)

			mockPhraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{
				ID:     1,
				UserID: 1,
				Phrase: "Test Phrase",
			}, nil)

			mockStorage.On("Upload", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)

			mockAudioRepo.On("Store", mock.Anything, mock.MatchedBy(func(audio *entity.Audio) bool {
				return audio.UserID == 1 && audio.PhraseID == 1 && audio.Status == "pending"
			})).Return(&entity.Audio{
				ID:            1,
				UserID:        1,
				PhraseID:      1,
				StoragePath:   "test-path",
				Status:        "pending",
				OriginalName:  "test.mp3",
				CurrentFormat: ".mp3",
			}, nil)

			mockQueue.On("Enqueue", mock.Anything, uint(1)).Return(nil)

			// Create use cases
			uploadUseCase := usecase.NewUploadAudioUseCase(
				mockAudioRepo,
				mockStorage,
				mockQueue,
				mockUserRepo,
				mockPhraseRepo,
			)

			downloadUseCase := usecase.NewDownloadAudioUseCase(
				mockAudioRepo,
				mockStorage,
				mockConverter,
				mockUserRepo,
				mockPhraseRepo,
			)

			// Create handler
			handler := handler.NewAudioHandler(uploadUseCase, downloadUseCase)

			// Create a multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("audio_file", "test.mp3")
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.Copy(part, strings.NewReader(tt.fileContent))
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()

			// Create request
			req := httptest.NewRequest("POST", "/users/"+tt.userID+"/phrases/"+tt.phraseID+"/audio", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Set up router with URL parameters
			router := mux.NewRouter()
			router.HandleFunc("/users/{user_id}/phrases/{phrase_id}/audio", handler.UploadAudio).Methods("POST")

			// Record response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response body
			if tt.expectedBody != nil {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if got["id"] != tt.expectedBody["id"] || got["status"] != tt.expectedBody["status"] {
					t.Errorf("handler returned unexpected body: got %v want %v", got, tt.expectedBody)
				}
			}

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
			mockPhraseRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
			mockAudioRepo.AssertExpectations(t)
			mockQueue.AssertExpectations(t)
		})
	}
}

func TestAudioHandler_Download(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		phraseID       string
		format         string
		expectedStatus int
	}{
		{
			name:           "successful download",
			userID:         "1",
			phraseID:       "1",
			format:         "mp3",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			mockAudioRepo := repoMocks.NewMockAudioRepository(t)
			mockUserRepo := repoMocks.NewMockUserRepository(t)
			mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)
			mockStorage := storageMocks.NewMockStorage(t)
			mockConverter := converterMocks.NewMockAudioConverter(t)
			mockQueue := queueMocks.NewMockTaskQueue(t)

			// Set up mock expectations
			mockUserRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.User{
				ID:   1,
				Name: "Test User",
			}, nil)

			mockPhraseRepo.On("GetByID", mock.Anything, uint(1)).Return(&entity.Phrase{
				ID:     1,
				UserID: 1,
				Phrase: "Test Phrase",
			}, nil)

			mockAudioRepo.On("GetByUserIDAndPhraseID", mock.Anything, uint(1), uint(1)).Return(&entity.Audio{
				ID:            1,
				UserID:        1,
				PhraseID:      1,
				StoragePath:   "test-path",
				Status:        "processed",
				OriginalName:  "test.mp3",
				CurrentFormat: ".mp3",
			}, nil)

			mockStorage.On("Download", mock.Anything, "test-path").Return(io.NopCloser(strings.NewReader("test audio content")), nil)

			mockConverter.On("ConvertFromReader", mock.Anything, mock.Anything, ".mp3", "mp3").Return(io.NopCloser(strings.NewReader("converted content")), nil)

			// Create use cases
			uploadUseCase := usecase.NewUploadAudioUseCase(
				mockAudioRepo,
				mockStorage,
				mockQueue,
				mockUserRepo,
				mockPhraseRepo,
			)

			downloadUseCase := usecase.NewDownloadAudioUseCase(
				mockAudioRepo,
				mockStorage,
				mockConverter,
				mockUserRepo,
				mockPhraseRepo,
			)

			// Create handler
			handler := handler.NewAudioHandler(uploadUseCase, downloadUseCase)

			// Create request
			req := httptest.NewRequest("GET", "/users/"+tt.userID+"/phrases/"+tt.phraseID+"/audio/"+tt.format, nil)

			// Set up router with URL parameters
			router := mux.NewRouter()
			router.HandleFunc("/users/{user_id}/phrases/{phrase_id}/audio/{format}", handler.GetAudio).Methods("GET")

			// Record response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// For successful responses, check that we got some content
			if tt.expectedStatus == http.StatusOK && rr.Body.Len() == 0 {
				t.Error("Response body is empty")
			}

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
			mockPhraseRepo.AssertExpectations(t)
			mockAudioRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
			mockConverter.AssertExpectations(t)
		})
	}
}
