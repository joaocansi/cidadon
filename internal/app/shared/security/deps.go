package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CryptoProvider interface {
	Hash(password string) string
	Compare(password, hash string) bool
}

type Jwt struct {
	Token     string
	ExpiresAt time.Time
}

type JwtProvider interface {
	Generate(sub string) (*Jwt, error)
	Verify(tokenString string) (*jwt.Token, error)
	GetSubject(tokenString string) (string, error)
}

type RefreshToken struct {
	Value     string
	ExpiresAt time.Time
	Hash      string
}

type RefreshTokenProvider interface {
	Generate() (*RefreshToken, error)
	Hash(token string) string
	Compare(token, hash string) bool
}
