package util

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemporaryFiles(t *testing.T) {
	tests := []struct {
		name           string
		originalFormat string
		targetFormat   string
	}{
		{
			name:           "m4a to mp3",
			originalFormat: "m4a",
			targetFormat:   "mp3",
		},
		{
			name:           "wav to flac",
			originalFormat: "wav",
			targetFormat:   "flac",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputPath, outputPath, err := CreateTemporaryFiles(tt.originalFormat, tt.targetFormat)
			assert.NoError(t, err)

			// Check that paths are in temp directory
			assert.Contains(t, inputPath, os.TempDir())
			assert.Contains(t, outputPath, os.TempDir())

			// Check file extensions
			assert.Equal(t, tt.originalFormat, path.Ext(inputPath)[1:])
			assert.Equal(t, tt.targetFormat, path.Ext(outputPath)[1:])

			// Check that paths are different
			assert.NotEqual(t, inputPath, outputPath)

			// Check that ULID parts are the same in both paths
			inputULID := path.Base(inputPath)[6:32] // Extract ULID part (after "audio_")
			outputULID := path.Base(outputPath)[6:32]
			assert.Equal(t, inputULID, outputULID)
		})
	}
}

func TestCleanupReadCloser(t *testing.T) {
	tests := []struct {
		name          string
		readerContent string
		cleanupErr    error
		closeErr      error
		expectedErr   error
	}{
		{
			name:          "successful cleanup and close",
			readerContent: "test content",
			cleanupErr:    nil,
			closeErr:      nil,
			expectedErr:   nil,
		},
		{
			name:          "cleanup error",
			readerContent: "test content",
			cleanupErr:    assert.AnError,
			closeErr:      nil,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupCalled := false

			mockReader := io.NopCloser(strings.NewReader(tt.readerContent))
			cleanup := func() error {
				cleanupCalled = true
				return tt.cleanupErr
			}

			crc := &CleanupReadCloser{
				ReadCloser: mockReader,
				Cleanup:    cleanup,
			}

			// Read content
			content, err := io.ReadAll(crc)
			assert.NoError(t, err)
			assert.Equal(t, tt.readerContent, string(content))

			// Close and check cleanup
			err = crc.Close()
			assert.Equal(t, tt.expectedErr, err)
			assert.True(t, cleanupCalled, "Cleanup function should have been called")
		})
	}
}
