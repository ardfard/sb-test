package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ardfard/sb-test/internal/usecase"
)

type AudioHandler struct {
	useCase *usecase.AudioUseCase
}

func NewAudioHandler(useCase *usecase.AudioUseCase) *AudioHandler {
	return &AudioHandler{
		useCase: useCase,
	}
}

func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	audio, err := h.useCase.UploadAudio(r.Context(), header.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with audio ID and status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     audio.ID,
		"status": audio.Status,
	})
}
