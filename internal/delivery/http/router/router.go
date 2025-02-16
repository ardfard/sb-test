package router

import (
	"net/http"

	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/gorilla/mux"
)

// SetupRoutes sets up all the HTTP routes for the application.
func SetupRoutes(audioHandler *handler.AudioHandler, userHandler *handler.UserHandler, phraseHandler *handler.PhraseHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/audio/user/{user_id}/phrase/{phrase_id}", audioHandler.UploadAudio).Methods(http.MethodPost)
	router.HandleFunc("/audio/user/{user_id}/phrase/{phrase_id}/{format}", audioHandler.GetAudio).Methods(http.MethodGet)

	router.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost)

	// Phrase routes
	router.HandleFunc("/users/{user_id}/phrases", phraseHandler.Create).Methods(http.MethodPost)

	return router
}
