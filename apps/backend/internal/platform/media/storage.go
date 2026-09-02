package media

import (
	"cidadon/internal/adapters/external/provider"
	environment "cidadon/internal/platform/config"
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewStorage(settings environment.Media) (provider.MediaStorage, error) {
	switch strings.ToLower(strings.TrimSpace(settings.Driver)) {
	case "", "local":
		return provider.NewLocalMediaStorage(settings.LocalDir, settings.PublicBaseURL)
	case "s3":
		if unusable(settings.S3Bucket) || unusable(settings.S3Region) {
			return nil, fmt.Errorf("MEDIA_S3_BUCKET and MEDIA_S3_REGION are required when MEDIA_DRIVER=s3")
		}
		accessKeyID := emptyWhenUnused(settings.S3AccessKeyID)
		secretAccessKey := emptyWhenUnused(settings.S3SecretAccessKey)
		endpoint := emptyWhenUnused(settings.S3Endpoint)
		if (accessKeyID == "") != (secretAccessKey == "") {
			return nil, fmt.Errorf("MEDIA_S3_ACCESS_KEY_ID and MEDIA_S3_SECRET_ACCESS_KEY must be configured together")
		}
		options := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(settings.S3Region)}
		if accessKeyID != "" || secretAccessKey != "" {
			options = append(options, awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")))
		}
		if endpoint != "" {
			options = append(options, awsconfig.WithBaseEndpoint(endpoint))
		}
		awsSettings, err := awsconfig.LoadDefaultConfig(context.Background(), options...)
		if err != nil {
			return nil, fmt.Errorf("load S3 configuration: %w", err)
		}
		return provider.NewS3MediaStorage(s3.NewFromConfig(awsSettings), settings.S3Bucket, settings.PublicBaseURL)
	default:
		return nil, fmt.Errorf("unsupported media driver %q", settings.Driver)
	}
}

func unusable(value string) bool { return strings.TrimSpace(value) == "" || value == "__unused__" }
func emptyWhenUnused(value string) string {
	if value == "__unused__" {
		return ""
	}
	return value
}
