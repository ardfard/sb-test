package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

func NewGCSStorage(client *storage.Client, bucketName string) *GCSStorage {
	return &GCSStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (g *GCSStorage) Upload(ctx context.Context, objectName string, reader io.Reader) error {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("failed to copy data to GCS: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}
