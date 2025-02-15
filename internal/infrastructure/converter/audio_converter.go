package converter

import (
	"context"
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type AudioConverter struct{}

func NewAudioConverter() *AudioConverter {
	return &AudioConverter{}
}

func (ac *AudioConverter) ConvertToWAV(ctx context.Context, inputPath, outputPath string) error {
	return ac.Convert(ctx, inputPath, outputPath, "wav")
}

func (ac *AudioConverter) Convert(ctx context.Context, inputPath, outputPath, outputFormat string) error {
	stream := ffmpeg.Input(inputPath)

	switch outputFormat {
	case "wav":
		stream = stream.Output(outputPath,
			ffmpeg.KwArgs{
				"acodec": "pcm_s16le",
				"ar":     "44100",
			})
	case "mp3":
		stream = stream.Output(outputPath,
			ffmpeg.KwArgs{
				"acodec": "libmp3lame",
				"q:a":    "2",
			})
	case "m4a":
		stream = stream.Output(outputPath,
			ffmpeg.KwArgs{
				"acodec": "aac",
				"strict": "experimental",
			})
	case "flac":
		stream = stream.Output(outputPath,
			ffmpeg.KwArgs{
				"acodec": "flac",
			})
	default:
		stream = stream.Output(outputPath)
	}

	err := stream.OverWriteOutput().Run()
	if err != nil {
		return fmt.Errorf("failed to convert audio: %v", err)
	}

	return nil
}
