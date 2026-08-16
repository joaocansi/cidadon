package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"time"
)

type RefreshTokenProviderImpl struct {
	expirationTime time.Duration
}

func NewRefreshTokenProvider(expirationTime time.Duration) *RefreshTokenProviderImpl {
	return &RefreshTokenProviderImpl{
		expirationTime: expirationTime,
	}
}

func (r *RefreshTokenProviderImpl) Generate() (*RefreshToken, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(token))

	tokenHash := base64.RawURLEncoding.EncodeToString(hash[:])
	return &RefreshToken{
		Value:     token,
		ExpiresAt: time.Now().Add(r.expirationTime),
		Hash:      tokenHash,
	}, nil
}

func (r *RefreshTokenProviderImpl) Hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (r *RefreshTokenProviderImpl) Compare(token, hash string) bool {
	return subtle.ConstantTimeCompare(
		[]byte(r.Hash(token)),
		[]byte(hash),
	) == 1
}
