package settings

import (
	"fmt"
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

type Settings struct {
	JwtProvider          JwtProvider
	RefreshTokenProvider RefreshTokenProvider
	Database             Database
}

var Env Settings

func Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := envconfig.InitWithOptions(&Env, envconfig.Options{AllOptional: false}); err != nil {
		return err
	}
	fmt.Println(Env)
	return nil
}
