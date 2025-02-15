package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Storage implements the Storage interface for AWS S3.
type S3Storage struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
}

// NewS3Storage creates a new S3Storage instance.
func NewS3Storage(region, bucket, accessKeyID, secretAccessKey string) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID, secretAccessKey, "",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %v", err)
	}
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)
	return &S3Storage{
		uploader:   uploader,
		downloader: downloader,
		bucket:     bucket,
	}, nil
}

// Upload uploads a file to S3.
func (s *S3Storage) Upload(ctx context.Context, objectName string, reader io.Reader) error {
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to s3: %v", err)
	}
	return nil
}

// Download downloads a file from S3.
func (s *S3Storage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	// Create a buffer to write our S3 object contents to.
	buf := aws.NewWriteAtBuffer([]byte{})

	// Create the input parameter for the download.
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectName),
	}

	// Download the file into the buffer.
	_, err := s.downloader.DownloadWithContext(ctx, buf, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from s3: %v", err)
	}

	// Return a ReadCloser wrapping a reader over the buffer's bytes.
	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
