package converter

import (
	"context"
	"io"
)

type AudioConverter interface {
	// Convert converts an audio file from inputPath to outputPath in the specified outputFormat
	Convert(ctx context.Context, inputPath, outputPath, outputFormat string) error

	// ConvertFromReader converts an audio file from input to output in the specified outputFormat
	ConvertFromReader(ctx context.Context, reader io.Reader, originalFormat, outputFormat string) (io.ReadCloser, error)
}
