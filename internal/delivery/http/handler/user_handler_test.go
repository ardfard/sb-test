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
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"name": "John Doe",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":   float64(1),
				"name": "John Doe",
			},
		},
		{
			name: "empty name",
			requestBody: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock
			mockUserRepo := repoMocks.NewMockUserRepository(t)

			// Set up mock expectations
			if tt.requestBody["name"] != "" {
				mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == "John Doe"
				})).Return(&entity.User{
					ID:   1,
					Name: "John Doe",
				}, nil)
			}

			// Create use case
			createUserUseCase := usecase.NewCreateUserUseCase(mockUserRepo)

			// Create handler
			handler := handler.NewUserHandler(createUserUseCase)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			rr := httptest.NewRecorder()
			handler.Create(rr, req)

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

				if got["id"] != tt.expectedBody["id"] || got["name"] != tt.expectedBody["name"] {
					t.Errorf("handler returned unexpected body: got %v want %v", got, tt.expectedBody)
				}
			}

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
		})
	}
}
