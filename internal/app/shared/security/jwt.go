package security

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProviderImpl struct {
	publicKey      *rsa.PublicKey
	privateKey     *rsa.PrivateKey
	issuer         string
	audience       []string
	expirationTime time.Duration
}

func NewJwtProvider(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, issuer string, audience []string, expirationTime time.Duration) *JwtProviderImpl {
	return &JwtProviderImpl{publicKey: publicKey, privateKey: privateKey, expirationTime: expirationTime, audience: audience, issuer: issuer}
}

func (j *JwtProviderImpl) Generate(sub string) (*Jwt, error) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expirationTime * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    j.issuer,
		Audience:  j.audience,
	}

	token, err := t.SignedString(j.privateKey)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		Token:     token,
		ExpiresAt: time.Now().Add(j.expirationTime),
	}, nil
}

func (j *JwtProviderImpl) Verify(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {

			// Prevent algorithm confusion attacks
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}

			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, errors.New("invalid signing algorithm")
			}

			return j.publicKey, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithAudience(j.audience...),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (j *JwtProviderImpl) GetSubject(tokenString string) (string, error) {
	token, err := j.Verify(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	if claims.Subject == "" {
		return "", errors.New("missing subject")
	}

	return claims.Subject, nil
}
