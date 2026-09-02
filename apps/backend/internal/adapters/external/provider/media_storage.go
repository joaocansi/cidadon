package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type StoredMedia struct {
	Key string
	URL string
}

// MediaStorage is the only outbound port used by the application for files.
// It keeps the rest of the application unaware of local disks or S3.
type MediaStorage interface {
	Store(ctx context.Context, prefix, contentType string, contents []byte) (StoredMedia, error)
	Delete(ctx context.Context, key string) error
}

func extensionForImage(contentType string) (string, error) {
	switch strings.ToLower(contentType) {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	default:
		return "", fmt.Errorf("unsupported media content type %q", contentType)
	}
}

func objectKey(prefix, id, extension string) string {
	prefix = strings.Trim(filepath.ToSlash(prefix), "/")
	if prefix == "" {
		prefix = "images"
	}
	return prefix + "/" + id + extension
}
