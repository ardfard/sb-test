package converter

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ardfard/sb-test/pkg/projectpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func TestAudioConverter(t *testing.T) {
	converter := NewAudioConverter()
	ctx := context.Background()
	tempDir := t.TempDir()

	testInputPath := filepath.Join(projectpath.RootProject, "tests", "fixtures", "test.m4a")

	t.Run("Convert to WAV", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output.wav")
		err := converter.Convert(ctx, testInputPath, outputPath, "wav")

		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		// Validate format
		probe, err := ffmpeg.Probe(outputPath)
		assert.NoError(t, err)
		assert.Contains(t, probe, "pcm_s16le")
		assert.Contains(t, probe, "44100")
	})

	t.Run("Convert to MP3", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output.mp3")
		err := converter.Convert(ctx, testInputPath, outputPath, "mp3")

		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		// Validate format
		probe, err := ffmpeg.Probe(outputPath)
		assert.NoError(t, err)
		assert.Contains(t, probe, "mp3")
	})

	t.Run("Convert to M4A", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output.m4a")
		err := converter.Convert(ctx, testInputPath, outputPath, "m4a")

		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		// Validate format
		probe, err := ffmpeg.Probe(outputPath)
		assert.NoError(t, err)
		assert.Contains(t, probe, "AAC")
	})

	t.Run("Convert to FLAC", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output.flac")
		err := converter.Convert(ctx, testInputPath, outputPath, "flac")

		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		// Validate format
		probe, err := ffmpeg.Probe(outputPath)
		assert.NoError(t, err)
		assert.Contains(t, probe, "FLAC")
	})

	t.Run("Convert with invalid input path", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output.mp3")
		err := converter.Convert(ctx, "nonexistent.mp3", outputPath, "mp3")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert audio")
	})

	t.Run("Convert with invalid output directory", func(t *testing.T) {
		outputPath := filepath.Join("/nonexistent/directory", "output.mp3")
		err := converter.Convert(ctx, testInputPath, outputPath, "mp3")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to convert audio")
	})
}

func TestAudioConverter_ConvertFromReader(t *testing.T) {
	tests := []struct {
		name          string
		inputFormat   string
		outputFormat  string
		expectedError bool
	}{
		{
			name:          "Success M4A to MP3",
			inputFormat:   "m4a",
			outputFormat:  "mp3",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testFile, err := os.Open(filepath.Join(projectpath.RootProject, "tests", "fixtures", "test.m4a"))
			require.NoError(t, err)
			defer testFile.Close()

			converter := NewAudioConverter()

			// Test
			result, err := converter.ConvertFromReader(context.Background(), testFile, tt.inputFormat, tt.outputFormat)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, result)
			defer result.Close()

			// Read the converted data
			data, err := io.ReadAll(result)
			assert.NoError(t, err)
			assert.NotEmpty(t, data)

			// Verify the converted data is in the correct format
			// This is a basic check - in a real scenario, you might want to use a library
			// to verify the audio format more thoroughly
			assert.Greater(t, len(data), 0)
		})
	}
}
