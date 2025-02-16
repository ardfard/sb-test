package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupRoutes(t *testing.T) {
	// Create mock handlers
	audioHandler := &handler.AudioHandler{}
	userHandler := &handler.UserHandler{}
	phraseHandler := &handler.PhraseHandler{}

	// Setup router
	router := SetupRoutes(audioHandler, userHandler, phraseHandler)
	require.NotNil(t, router)

	// Test cases for route configuration
	tests := []struct {
		name          string
		method        string
		path          string
		expectedRoute bool
	}{
		{
			name:          "Audio Upload Route",
			method:        http.MethodPost,
			path:          "/audio/user/1/phrase/1",
			expectedRoute: true,
		},
		{
			name:          "Audio Get Route",
			method:        http.MethodGet,
			path:          "/audio/user/1/phrase/1/mp3",
			expectedRoute: true,
		},
		{
			name:          "User Create Route",
			method:        http.MethodPost,
			path:          "/users",
			expectedRoute: true,
		},
		{
			name:          "Phrase Create Route",
			method:        http.MethodPost,
			path:          "/users/1/phrases",
			expectedRoute: true,
		},
		{
			name:          "Health Check Route",
			method:        http.MethodGet,
			path:          "/health",
			expectedRoute: true,
		},
		{
			name:          "Non-existent Route",
			method:        http.MethodGet,
			path:          "/nonexistent",
			expectedRoute: false,
		},
		{
			name:          "Wrong Method for Existing Path",
			method:        http.MethodDelete,
			path:          "/users",
			expectedRoute: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request to test route matching
			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			// Test if route exists
			var routeMatch mux.RouteMatch
			matched := router.Match(req, &routeMatch)

			// Assert route matching
			assert.Equal(t, tt.expectedRoute, matched, "Route matching failed for %s %s", tt.method, tt.path)

			if matched {
				// Additional checks for matched routes
				assert.NotNil(t, routeMatch.Handler, "Handler should not be nil for matched route")
				assert.NotNil(t, routeMatch.Route, "Route should not be nil for matched route")
			}
		})
	}
}

func TestSetupRoutes_HandlerNilCheck(t *testing.T) {
	// Test with nil handlers
	router := SetupRoutes(nil, nil, nil)
	require.NotNil(t, router, "Router should be created even with nil handlers")

	// Basic check that health endpoint still works
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	var routeMatch mux.RouteMatch
	matched := router.Match(req, &routeMatch)
	assert.True(t, matched, "Health route should still be accessible")
}

func TestRouteHandlers(t *testing.T) {
	// Create mock handlers
	audioHandler := &handler.AudioHandler{}
	userHandler := &handler.UserHandler{}
	phraseHandler := &handler.PhraseHandler{}

	// Setup router
	router := SetupRoutes(audioHandler, userHandler, phraseHandler)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test cases for actual HTTP requests
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Health Check",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found Route",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and execute request
			url := server.URL + tt.path
			req, err := http.NewRequest(tt.method, url, nil)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
