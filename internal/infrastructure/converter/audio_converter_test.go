package converter

import (
	"context"
	"os"
	"testing"
)

func TestConvertToWAVWithNonExistentFile(t *testing.T) {
	ac := NewAudioConverter()
	ctx := context.Background()
	inputPath := "nonexistent_input.wav"
	outputPath := "nonexistent_output.wav"

	err := ac.ConvertToWAV(ctx, inputPath, outputPath)
	if err == nil {
		t.Fatal("expected error when converting nonexistent file, got nil")
	}
}

func TestConvertToWAV_Success(t *testing.T) {
	// This test is optional; it attempts to convert a sample file if available.
	// If the sample file "sample.mp3" is not found, the test is skipped.
	sampleInput := "sample.mp3"
	outputPath := "test_output.wav"
	if _, err := os.Stat(sampleInput); os.IsNotExist(err) {
		t.Skip("sample.mp3 not found, skipping conversion success test")
	}

	ac := NewAudioConverter()
	ctx := context.Background()
	err := ac.ConvertToWAV(ctx, sampleInput, outputPath)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}
	// Clean up output file
	os.Remove(outputPath)
}
