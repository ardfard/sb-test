package converter

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ardfard/sb-test/pkg/util"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type AudioConverter struct{}

func NewAudioConverter() *AudioConverter {
	return &AudioConverter{}
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

func (ac *AudioConverter) ConvertFromReader(ctx context.Context, input io.Reader, originalFormat, outputFormat string) (io.ReadCloser, error) {

	inputPath, outputPath, err := util.CreateTemporaryFiles(originalFormat, outputFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary files: %v", err)
	}
	defer os.Remove(inputPath)

	inputFile, err := os.Create(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer inputFile.Close()

	if _, err := io.Copy(inputFile, input); err != nil {
		return nil, fmt.Errorf("failed to write to temp input file: %v", err)
	}

	err = ac.Convert(ctx, inputPath, outputPath, outputFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to convert audio: %v", err)
	}

	outputFile, err := os.Open(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %v", err)
	}

	return &util.CleanupReadCloser{
		ReadCloser: outputFile,
		Cleanup: func() error {
			return os.Remove(outputPath)
		},
	}, nil
}
