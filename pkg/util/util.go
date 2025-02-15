package util

import (
	"fmt"

	"github.com/ardfard/sb-test/internal/domain/entity"
)

func CreateTemporaryFiles(audio *entity.Audio, targetFormat string) (string, string, error) {
	inputPath := fmt.Sprintf("/tmp/%d%s", audio.ID, audio.OriginalFormat)
	outputPath := fmt.Sprintf("/tmp/%d.%s", audio.ID, targetFormat)
	return inputPath, outputPath, nil
}
