package storage

import (
	"fmt"

	"github.com/ardfard/sb-test/config"
	"github.com/ardfard/sb-test/internal/domain/storage"
)

// NewStorage creates a new storage instance based on the configuration.
func NewStorage(cfg *config.StorageConfig) (storageInstance storage.Storage, err error) {
	if cfg == nil {
		return nil, fmt.Errorf("storage configuration is nil")
	}
	switch cfg.Type {
	// Initialize GCS client
	case "s3":
		storageInstance, err = NewS3Storage(
			cfg.S3.Region,
			cfg.S3.Bucket,
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
		)
		if err != nil {
			return nil, err
		}
	case "local":
		storageInstance, err = NewLocalStorage(cfg.Local.Directory)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
	return
}
