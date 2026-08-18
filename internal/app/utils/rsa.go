package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type RSA256Pair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func LoadRSAKeys(privatePath, publicPath string) (*RSA256Pair, error) {
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
