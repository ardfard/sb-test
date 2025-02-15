package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements the Storage interface using the local filesystem.
type LocalStorage struct {
	directory string
}

// NewLocalStorage creates a new LocalStorage instance.
func NewLocalStorage(directory string) (*LocalStorage, error) {
	// Ensure the directory exists.
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create local storage directory: %v", err)
	}
	return &LocalStorage{
		directory: directory,
	}, nil
}

// Upload saves the file to the local filesystem.
func (ls *LocalStorage) Upload(ctx context.Context, objectName string, reader io.Reader) error {
	filePath := filepath.Join(ls.directory, objectName)
	// Create any missing directories.
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for file: %v", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

// Download downloads the file from the local filesystem.
func (ls *LocalStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	filePath := filepath.Join(ls.directory, objectName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return file, nil
}
