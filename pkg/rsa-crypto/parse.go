package rsacrypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

func ParsePublicKey(publicKeyData []byte) (*rsa.PublicKey, error) {
	pemBlock, _ := pem.Decode(publicKeyData)
	if pemBlock == nil {
		return nil, errors.New("no PEM block found")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("unsupported type of public key")
	}

	return rsaPubKey, nil
}

func ParsePrivateKey(privateKeyData []byte) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode(privateKeyData)
	if pemBlock == nil {
		return nil, errors.New("no PEM block found")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %w", err)
	}

	return privateKey, nil
}

func FromFile[Key rsa.PrivateKey | rsa.PublicKey](path string, parseFn func([]byte) (*Key, error)) (*Key, error) {
	keyFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file containing public key: %w", err)
	}
	defer keyFile.Close()

	keyBytes, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file containing public key: %w", err)
	}

	key, err := parseFn(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file contebt into key: %w", err)
	}

	return key, nil
}
