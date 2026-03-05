package ports

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, filename string, contentType string, size int64, data io.Reader) (string, error)
	Delete(ctx context.Context, objectName string) error
}
