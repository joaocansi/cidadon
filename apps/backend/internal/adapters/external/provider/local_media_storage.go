package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LocalMediaStorage struct {
	rootDir   string
	publicURL string
}

func NewLocalMediaStorage(rootDir, publicURL string) (*LocalMediaStorage, error) {
	if strings.TrimSpace(rootDir) == "" {
		return nil, fmt.Errorf("local media directory is required")
	}
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return nil, fmt.Errorf("create local media directory: %w", err)
	}
	return &LocalMediaStorage{rootDir: rootDir, publicURL: strings.TrimRight(publicURL, "/")}, nil
}

func (s *LocalMediaStorage) Store(_ context.Context, prefix, contentType string, contents []byte) (StoredMedia, error) {
	extension, err := extensionForImage(contentType)
	if err != nil {
		return StoredMedia{}, err
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return StoredMedia{}, fmt.Errorf("generate media key: %w", err)
	}
	key := objectKey(prefix, hex.EncodeToString(random[:]), extension)
	path := filepath.Join(s.rootDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return StoredMedia{}, fmt.Errorf("create media directory: %w", err)
	}
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		return StoredMedia{}, fmt.Errorf("write media file: %w", err)
	}
	return StoredMedia{Key: key, URL: s.publicURL + "/" + key}, nil
}

func (s *LocalMediaStorage) Delete(_ context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return nil
	}
	if err := os.Remove(filepath.Join(s.rootDir, filepath.FromSlash(key))); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete local media: %w", err)
	}
	return nil
}
