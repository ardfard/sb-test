package converter

import (
	"context"
	"fmt"
	"os/exec"
)

type AudioConverter struct{}

func NewAudioConverter() *AudioConverter {
	return &AudioConverter{}
}

func (ac *AudioConverter) ConvertToWAV(ctx context.Context, inputPath, outputPath string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inputPath, "-acodec", "pcm_s16le", "-ar", "44100", outputPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert audio: %v", err)
	}

	return nil
}
