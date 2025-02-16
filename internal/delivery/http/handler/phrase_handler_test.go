package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/ardfard/sb-test/internal/domain/entity"
	repoMocks "github.com/ardfard/sb-test/internal/infrastructure/repository/mocks"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestPhraseHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful creation",
			method: "POST",
			userID: "1",
			requestBody: map[string]interface{}{
				"text": "Hello, World!",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":   float64(1),
				"text": "Hello, World!",
			},
		},
		{
			name:   "empty text",
			method: "POST",
			userID: "1",
			requestBody: map[string]interface{}{
				"text": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid user ID",
			method: "POST",
			userID: "invalid",
			requestBody: map[string]interface{}{
				"text": "Hello, World!",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "method not allowed",
			method:         "GET",
			userID:         "1",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid request body",
			method:         "POST",
			userID:         "1",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			mockPhraseRepo := repoMocks.NewMockPhraseRepository(t)

			// Use mock.MatchedBy for more flexible matching
			if tt.requestBody != nil && tt.method == "POST" && tt.requestBody["text"] != "" && tt.userID == "1" {
				mockPhraseRepo.On("Create", mock.Anything, mock.MatchedBy(func(phrase *entity.Phrase) bool {
					return phrase.UserID == 1 && phrase.Phrase == "Hello, World!"
				})).Return(&entity.Phrase{
					ID:     1,
					UserID: 1,
					Phrase: "Hello, World!",
				}, nil)
			}

			// Create use case
			createPhraseUseCase := usecase.NewCreatePhraseUseCase(mockPhraseRepo)

			// Create handler
			handler := handler.NewPhraseHandler(createPhraseUseCase)

			// Create request
			var req *http.Request
			if tt.requestBody != nil {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(tt.method, "/users/"+tt.userID+"/phrases", bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(tt.method, "/users/"+tt.userID+"/phrases", bytes.NewBuffer([]byte("invalid json")))
			}
			req.Header.Set("Content-Type", "application/json")

			// Set up router with URL parameters
			router := mux.NewRouter()
			router.HandleFunc("/users/{user_id}/phrases", handler.Create).Methods("POST")

			// Record response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response body for successful cases
			if tt.expectedStatus == http.StatusOK {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if got["id"] != tt.expectedBody["id"] || got["text"] != tt.expectedBody["text"] {
					t.Errorf("handler returned unexpected body: got %v want %v", got, tt.expectedBody)
				}
			}

			// Verify all expectations were met
			mockPhraseRepo.AssertExpectations(t)
		})
	}
}
