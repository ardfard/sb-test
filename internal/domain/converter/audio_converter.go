package converter

import "context"

type AudioConverter interface {
	Convert(ctx context.Context, inputPath, outputPath, outputFormat string) error
}
