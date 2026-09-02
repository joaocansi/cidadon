package environment

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JwtProvider struct {
	Issuer     string
	Expiration time.Duration
	Audience   []string
	PrivateKey string
	PublicKey  string
}

type RefreshTokenProvider struct {
	Expiration time.Duration
}

type Media struct {
	Driver            string
	LocalDir          string
	PublicBaseURL     string
	S3Bucket          string
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
}

type Settings struct {
	JwtProvider          JwtProvider
	RefreshTokenProvider RefreshTokenProvider
	Database             Database
	Media                Media
}

var Env Settings

func Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	setMediaDefaults()
	if err := envconfig.InitWithOptions(&Env, envconfig.Options{AllOptional: false}); err != nil {
		return err
	}
	return nil
}

func setMediaDefaults() {
	defaults := map[string]string{
		"MEDIA_DRIVER":               "local",
		"MEDIA_LOCAL_DIR":            "../../.runtime/uploads",
		"MEDIA_PUBLIC_BASE_URL":      "http://localhost:8080/media",
		"MEDIA_S3_BUCKET":            "__unused__",
		"MEDIA_S3_REGION":            "__unused__",
		"MEDIA_S3_ENDPOINT":          "__unused__",
		"MEDIA_S3_ACCESS_KEY_ID":     "__unused__",
		"MEDIA_S3_SECRET_ACCESS_KEY": "__unused__",
	}
	for key, value := range defaults {
		if current, found := os.LookupEnv(key); !found || current == "" {
			_ = os.Setenv(key, value)
		}
	}
}

type RSA256Pair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func LoadRsaKeys(privatePath, publicPath string) (*RSA256Pair, error) {
	privateBytes, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, err
	}

	privateBlock, _ := pem.Decode(privateBytes)

	key, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	publicBytes, err := os.ReadFile(publicPath)
	if err != nil {
		return nil, err
	}

	publicBlock, _ := pem.Decode(publicBytes)

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)

	return &RSA256Pair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
