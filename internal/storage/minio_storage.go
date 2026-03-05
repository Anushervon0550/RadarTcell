package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName, publicURL string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client init: %w", err)
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
		publicURL:  strings.TrimRight(publicURL, "/"),
	}, nil
}

func (s *MinioStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("check bucket exists: %w", err)
	}

	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		// Set public read policy
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, s.bucketName)

		if err := s.client.SetBucketPolicy(ctx, s.bucketName, policy); err != nil {
			return fmt.Errorf("set bucket policy: %w", err)
		}
	}

	return nil
}

func (s *MinioStorage) Upload(ctx context.Context, filename string, contentType string, size int64, data io.Reader) (string, error) {
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload object: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, objectName)
	return url, nil
}

func (s *MinioStorage) Delete(ctx context.Context, objectName string) error {
	if err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

