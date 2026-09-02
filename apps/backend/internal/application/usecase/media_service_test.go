package usecase

import (
	"cidadon/internal/adapters/external/provider"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const onePixelPNG = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScLJ4QAAAABJRU5ErkJggg=="

func TestStoreDataURLWritesAValidatedImage(t *testing.T) {
	directory := t.TempDir()
	storage, err := provider.NewLocalMediaStorage(directory, "http://localhost:8080/media")
	if err != nil {
		t.Fatal(err)
	}
	service := NewMediaService(storage)
	item, err := service.StoreDataURL(context.Background(), "demands", onePixelPNG)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(item.URL, "http://localhost:8080/media/demands/") {
		t.Fatalf("unexpected URL %q", item.URL)
	}
	if _, err := os.Stat(filepath.Join(directory, filepath.FromSlash(item.Key))); err != nil {
		t.Fatalf("stored file does not exist: %v", err)
	}
}

func TestStoreDataURLRejectsUnsupportedImage(t *testing.T) {
	storage, err := provider.NewLocalMediaStorage(t.TempDir(), "http://localhost:8080/media")
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewMediaService(storage).StoreDataURL(context.Background(), "demands", "data:image/gif;base64,R0lGODlhAQABAAAAACw=")
	if err == nil {
		t.Fatal("expected GIF to be rejected")
	}
}
