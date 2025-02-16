package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/ardfard/sb-test/pkg/logger"
)

type UserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
}

func NewUserHandler(createUserUseCase *usecase.CreateUserUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
	}
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		logger.Error("Name is required")
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	user, err := h.createUserUseCase.Create(r.Context(), req.Name)
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := CreateUserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
