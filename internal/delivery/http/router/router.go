package router

import (
	"net/http"

	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/gorilla/mux" // Alternatively, you can use the standard net/http
)

// SetupRoutes sets up all the HTTP routes for the application.
func SetupRoutes(audioHandler *handler.AudioHandler) *mux.Router {
	router := mux.NewRouter()

	// Define your routes here:
	router.HandleFunc("/upload", audioHandler.UploadAudio).Methods(http.MethodPost)
	// Add more routes as the application evolves.

	return router
}
