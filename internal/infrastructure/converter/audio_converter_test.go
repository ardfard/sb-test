package converter

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ardfard/sb-test/pkg/projectpath"
	"github.com/stretchr/testify/assert"
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
