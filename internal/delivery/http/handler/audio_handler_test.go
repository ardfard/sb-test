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
		method         string
		userID         string
		phraseID       string
		fileContent    string
		skipFile       bool
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful upload",
			method:         "POST",
			userID:         "1",
			phraseID:       "1",
			fileContent:    "test audio content",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":     float64(1),
				"status": "pending",
			},
		},
		{
			name:           "method not allowed",
			method:         "GET",
			userID:         "1",
			phraseID:       "1",
			fileContent:    "test audio content",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid user ID",
			method:         "POST",
			userID:         "invalid",
			phraseID:       "1",
			fileContent:    "test audio content",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid phrase ID",
			method:         "POST",
			userID:         "1",
			phraseID:       "invalid",
			fileContent:    "test audio content",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing file",
			method:         "POST",
			userID:         "1",
			phraseID:       "1",
			skipFile:       true,
			expectedStatus: http.StatusBadRequest,
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

			// Set up mock expectations only for successful case
			if tt.expectedStatus == http.StatusOK {
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
			}

			// Create use cases and handler
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

			handler := handler.NewAudioHandler(uploadUseCase, downloadUseCase)

			var req *http.Request
			if !tt.skipFile {
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

				req = httptest.NewRequest(tt.method, "/users/"+tt.userID+"/phrases/"+tt.phraseID+"/audio", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				req = httptest.NewRequest(tt.method, "/users/"+tt.userID+"/phrases/"+tt.phraseID+"/audio", nil)
			}

			router := mux.NewRouter()
			router.HandleFunc("/users/{user_id}/phrases/{phrase_id}/audio", handler.UploadAudio).Methods("POST")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

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
		method         string
		userID         string
		phraseID       string
		format         string
		expectedStatus int
	}{
		{
			name:           "successful download",
			method:         "GET",
			userID:         "1",
			phraseID:       "1",
			format:         "mp3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "method not allowed",
			method:         "POST",
			userID:         "1",
			phraseID:       "1",
			format:         "mp3",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid user ID",
			method:         "GET",
			userID:         "invalid",
			phraseID:       "1",
			format:         "mp3",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid phrase ID",
			method:         "GET",
			userID:         "1",
			phraseID:       "invalid",
			format:         "mp3",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudioRepo := repoMocks.NewMockAudioRepository(t)
			mockUserRepo := repoMocks.NewMockUserRepository(t)
			mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)
			mockStorage := storageMocks.NewMockStorage(t)
			mockConverter := converterMocks.NewMockAudioConverter(t)
			mockQueue := queueMocks.NewMockTaskQueue(t)

			// Set up mock expectations only for successful case
			if tt.expectedStatus == http.StatusOK {
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
			}

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

			handler := handler.NewAudioHandler(uploadUseCase, downloadUseCase)

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID+"/phrases/"+tt.phraseID+"/audio/"+tt.format, nil)

			router := mux.NewRouter()
			router.HandleFunc("/users/{user_id}/phrases/{phrase_id}/audio/{format}", handler.GetAudio).Methods("GET")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK && rr.Body.Len() == 0 {
				t.Error("Response body is empty")
			}

			mockUserRepo.AssertExpectations(t)
			mockPhraseRepo.AssertExpectations(t)
			mockAudioRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
			mockConverter.AssertExpectations(t)
		})
	}
}
