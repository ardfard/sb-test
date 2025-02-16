package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/ardfard/sb-test/pkg/logger"
	"github.com/gorilla/mux"
)

type PhraseHandler struct {
	createPhraseUseCase *usecase.CreatePhraseUseCase
}

func NewPhraseHandler(createPhraseUseCase *usecase.CreatePhraseUseCase) *PhraseHandler {
	return &PhraseHandler{
		createPhraseUseCase: createPhraseUseCase,
	}
}

type CreatePhraseRequest struct {
	Text string `json:"text"`
}

type CreatePhraseResponse struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

func (h *PhraseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := mux.Vars(r)["user_id"]

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		logger.Errorf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req CreatePhraseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	phrase, err := h.createPhraseUseCase.Create(r.Context(), req.Text, uint(userIDUint))
	if err != nil {
		logger.Errorf("Failed to create phrase: %v", err)
		http.Error(w, "Failed to create phrase", http.StatusInternalServerError)
		return
	}

	response := CreatePhraseResponse{
		ID:   phrase.ID,
		Text: phrase.Phrase,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
