package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ardfard/sb-test/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.StorageConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
		{
			name: "invalid storage type",
			cfg: &config.StorageConfig{
				Type: "invalid",
			},
			wantErr: true,
		},
		{
			name: "valid local storage",
			cfg: &config.StorageConfig{
				Type: "local",
				Local: &config.LocalStorageConfig{
					Directory: t.TempDir(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewStorage(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, storage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, storage)
			}
		})
	}
}

func TestLocalStorage(t *testing.T) {
	tempDir := t.TempDir()
	localStorage, err := NewLocalStorage(tempDir)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("test content")
	objectName := "test/file.txt"

	// Test Upload
	t.Run("upload file", func(t *testing.T) {
		err := localStorage.Upload(ctx, objectName, bytes.NewReader(content))
		require.NoError(t, err)

		// Verify file exists
		filePath := filepath.Join(tempDir, objectName)
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	})

	// Test Download
	t.Run("download file", func(t *testing.T) {
		reader, err := localStorage.Download(ctx, objectName)
		require.NoError(t, err)
		defer reader.Close()

		downloaded, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, downloaded)
	})

	// Test Download non-existent file
	t.Run("download non-existent file", func(t *testing.T) {
		reader, err := localStorage.Download(ctx, "non-existent.txt")
		assert.Error(t, err)
		assert.Nil(t, reader)
	})
}

func TestS3Storage(t *testing.T) {
	// MinIO test setup
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	region := "us-east-1"
	bucket := "testbucket"
	useSSL := false

	// Skip if MinIO is not available
	if !isMinIOAvailable(endpoint) {
		t.Skip("MinIO is not available")
	}

	// Create test bucket
	sess := createMinIOSession(endpoint, accessKeyID, secretAccessKey, region, useSSL)
	s3Client := s3.New(sess)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	require.NoError(t, err)

	// Initialize S3Storage
	s3Storage, err := NewS3Storage(region, bucket, accessKeyID, secretAccessKey, true)
	require.NoError(t, err)

	ctx := context.Background()
	content := []byte("test content")
	objectName := "test/file.txt"

	// Test Upload
	t.Run("upload file", func(t *testing.T) {
		err := s3Storage.Upload(ctx, objectName, bytes.NewReader(content))
		require.NoError(t, err)
	})

	// Test Download
	t.Run("download file", func(t *testing.T) {
		reader, err := s3Storage.Download(ctx, objectName)
		require.NoError(t, err)
		defer reader.Close()

		downloaded, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, downloaded)
	})

	// Test Download non-existent file
	t.Run("download non-existent file", func(t *testing.T) {
		reader, err := s3Storage.Download(ctx, "non-existent.txt")
		assert.Error(t, err)
		assert.Nil(t, reader)
	})

	// Cleanup
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectName),
	})
	require.NoError(t, err)

	_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	require.NoError(t, err)
}

func isMinIOAvailable(endpoint string) bool {
	sess := createMinIOSession(endpoint, "minioadmin", "minioadmin", "us-east-1", false)
	s3Client := s3.New(sess)
	_, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	return err == nil
}

func createMinIOSession(endpoint, accessKeyID, secretAccessKey, region string, useSSL bool) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!useSSL),
		S3ForcePathStyle: aws.Bool(true),
	})
	return sess
}
