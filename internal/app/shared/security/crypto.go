package security

import (
	"golang.org/x/crypto/bcrypt"
)

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
