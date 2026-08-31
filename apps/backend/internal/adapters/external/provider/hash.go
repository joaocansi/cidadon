package provider

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

type Hash struct {
	Value string
	Hash  string
}

type HashProvider interface {
	Generate() (*Hash, error)
	Hash(value string) string
	Compare(value, hash string) bool
}

type HashProviderImpl struct{}

func NewHashProvider() *HashProviderImpl {
	return &HashProviderImpl{}
}

func (r *HashProviderImpl) Generate() (*Hash, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	value := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(value))

	valueHash := base64.RawURLEncoding.EncodeToString(hash[:])
	return &Hash{
		Value: value,
		Hash:  valueHash,
	}, nil
}

func (r *HashProviderImpl) Hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (r *HashProviderImpl) Compare(token, hash string) bool {
	return subtle.ConstantTimeCompare(
		[]byte(r.Hash(token)),
		[]byte(hash),
	) == 1
}
