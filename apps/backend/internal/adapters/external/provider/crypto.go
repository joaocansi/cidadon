package provider

import (
	"golang.org/x/crypto/bcrypt"
)

type CryptoProvider interface {
	Hash(password string) string
	Compare(password, hash string) bool
}

type CryptoProviderImpl struct {
	cost int
}

func NewCryptoProvider() CryptoProvider {
	return &CryptoProviderImpl{
		cost: bcrypt.DefaultCost,
	}
}

func (cp *CryptoProviderImpl) Hash(password string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		cp.cost,
	)

	if err != nil {
		panic(err)
	}

	return string(hash)
}

func (cp *CryptoProviderImpl) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
