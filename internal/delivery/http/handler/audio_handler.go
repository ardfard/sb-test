package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/gorilla/mux"
)

// AudioHandler is a handler for audio-related operations.
type AudioHandler struct {
	uploadUseCase   *usecase.UploadAudioUseCase
	downloadUseCase *usecase.DownloadAudioUseCase
}

// NewAudioHandler creates a new AudioHandler with the given use case.
func NewAudioHandler(uploadUseCase *usecase.UploadAudioUseCase, downloadUseCase *usecase.DownloadAudioUseCase) *AudioHandler {
	return &AudioHandler{
		uploadUseCase:   uploadUseCase,
		downloadUseCase: downloadUseCase,
	}
}

// UploadAudio handles the upload of an audio file.
func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := mux.Vars(r)["user_id"]
	phraseID := mux.Vars(r)["phrase_id"]

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	phraseIDUint, err := strconv.ParseUint(phraseID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid phrase ID", http.StatusBadRequest)
		return
	}

	audio, err := h.uploadUseCase.Upload(r.Context(), header.Filename, file, uint(userIDUint), uint(phraseIDUint))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with audio ID and status
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     audio.ID,
		"status": audio.Status,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := mux.Vars(r)["user_id"]
	phraseID := mux.Vars(r)["phrase_id"]
	format := mux.Vars(r)["format"]

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	phraseIDUint, err := strconv.ParseUint(phraseID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid phrase ID", http.StatusBadRequest)
		return
	}

	audio, err := h.downloadUseCase.Download(r.Context(), uint(userIDUint), uint(phraseIDUint), format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, audio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
