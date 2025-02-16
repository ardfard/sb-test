package handler

import (
	"net/http"

	"github.com/ardfard/sb-test/pkg/logger"
)

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"healthy"}`)); err != nil {
		logger.Errorf("failed to write response: %v", err)
	}
}
