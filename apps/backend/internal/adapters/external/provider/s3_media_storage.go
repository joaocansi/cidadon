package provider

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3MediaStorage struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewS3MediaStorage(client *s3.Client, bucket, publicURL string) (*S3MediaStorage, error) {
	if client == nil || strings.TrimSpace(bucket) == "" || strings.TrimSpace(publicURL) == "" {
		return nil, fmt.Errorf("S3 client, bucket and public URL are required")
	}
	return &S3MediaStorage{client: client, bucket: bucket, publicURL: strings.TrimRight(publicURL, "/")}, nil
}

func (s *S3MediaStorage) Store(ctx context.Context, prefix, contentType string, contents []byte) (StoredMedia, error) {
	extension, err := extensionForImage(contentType)
	if err != nil {
		return StoredMedia{}, err
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return StoredMedia{}, fmt.Errorf("generate media key: %w", err)
	}
	key := objectKey(prefix, hex.EncodeToString(random[:]), extension)
	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(contents),
		ContentType: aws.String(contentType),
	}); err != nil {
		return StoredMedia{}, fmt.Errorf("put S3 media: %w", err)
	}
	return StoredMedia{Key: key, URL: s.publicURL + "/" + key}, nil
}

func (s *S3MediaStorage) Delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return fmt.Errorf("delete S3 media: %w", err)
	}
	return nil
}
