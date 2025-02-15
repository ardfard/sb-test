package storage

import (
	"context"
	"io"
)

// Storage is a contract for any object storage provider.
type Storage interface {
	// Upload saves the data under a specified object name.
	Upload(ctx context.Context, objectName string, reader io.Reader) error
	// Download downloads the data from the specified object name.
	Download(ctx context.Context, objectName string) (io.ReadCloser, error)
}
