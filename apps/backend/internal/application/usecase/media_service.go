package usecase

import (
	"cidadon/internal/adapters/external/provider"
	service "cidadon/internal/application/contract"
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	maxImageBytes = 2 << 20
	maxImageCount = 5
)

type MediaService struct{ storage provider.MediaStorage }

func NewMediaService(storage provider.MediaStorage) *MediaService {
	return &MediaService{storage: storage}
}

func (s *MediaService) StoreFiles(ctx context.Context, prefix string, files []*multipart.FileHeader, maxCount int) ([]provider.StoredMedia, error) {
	if maxCount <= 0 {
		maxCount = maxImageCount
	}
	if len(files) > maxCount {
		return nil, service.InvalidInput("too many images").WithDetails(map[string]any{"fields": []string{"images"}})
	}
	stored := make([]provider.StoredMedia, 0, len(files))
	for _, file := range files {
		contents, contentType, err := readImage(file)
		if err != nil {
			s.DeleteAll(ctx, stored)
			return nil, err
		}
		item, err := s.storage.Store(ctx, prefix, contentType, contents)
		if err != nil {
			s.DeleteAll(ctx, stored)
			return nil, service.Unavailable("media storage unavailable")
		}
		stored = append(stored, item)
	}
	return stored, nil
}

func (s *MediaService) StoreDataURL(ctx context.Context, prefix, value string) (provider.StoredMedia, error) {
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "data:") || !strings.Contains(parts[0], ";base64") {
		return provider.StoredMedia{}, service.InvalidInput("invalid legacy image")
	}
	contents, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return provider.StoredMedia{}, service.InvalidInput("invalid legacy image")
	}
	contentType := http.DetectContentType(contents)
	if err := validateImage(contents, contentType); err != nil {
		return provider.StoredMedia{}, err
	}
	item, err := s.storage.Store(ctx, prefix, contentType, contents)
	if err != nil {
		return provider.StoredMedia{}, service.Unavailable("media storage unavailable")
	}
	return item, nil
}

func (s *MediaService) DeleteAll(ctx context.Context, items []provider.StoredMedia) {
	for _, item := range items {
		_ = s.storage.Delete(ctx, item.Key)
	}
}

func MediaURLs(items []provider.StoredMedia) []string {
	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.URL)
	}
	return urls
}

func readImage(file *multipart.FileHeader) ([]byte, string, error) {
	if file == nil {
		return nil, "", service.InvalidInput("image is required").WithDetails(map[string]any{"fields": []string{"photo"}})
	}
	if file.Size > maxImageBytes {
		return nil, "", service.InvalidInput("image is too large").WithDetails(map[string]any{"fields": []string{"images"}})
	}
	opened, err := file.Open()
	if err != nil {
		return nil, "", service.InvalidInput("invalid image")
	}
	defer opened.Close()
	contents, err := io.ReadAll(io.LimitReader(opened, maxImageBytes+1))
	if err != nil {
		return nil, "", service.InvalidInput("invalid image")
	}
	contentType := http.DetectContentType(contents)
	if err := validateImage(contents, contentType); err != nil {
		return nil, "", err
	}
	return contents, contentType, nil
}

func validateImage(contents []byte, contentType string) error {
	if len(contents) == 0 || len(contents) > maxImageBytes {
		return service.InvalidInput("invalid image size").WithDetails(map[string]any{"fields": []string{"images"}})
	}
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return service.InvalidInput("unsupported image type %q", contentType).WithDetails(map[string]any{"fields": []string{"images"}})
	}
}
