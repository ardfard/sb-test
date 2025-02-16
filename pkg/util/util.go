package util

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/oklog/ulid/v2"
)

func CreateTemporaryFiles(originalFormat, targetFormat string) (string, string, error) {
	ulid := ulid.Make()
	inputPath := path.Join(os.TempDir(), fmt.Sprintf("audio_%s.%s", ulid, originalFormat))
	outputPath := path.Join(os.TempDir(), fmt.Sprintf("audio_%s.%s", ulid, targetFormat))
	return inputPath, outputPath, nil
}

type CleanupReadCloser struct {
	io.ReadCloser
	Cleanup func() error
}

func (c *CleanupReadCloser) Close() error {
	err := c.ReadCloser.Close()
	if cleanupErr := c.Cleanup(); cleanupErr != nil {
		return cleanupErr
	}
	return err
}
