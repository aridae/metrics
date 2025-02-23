package rsacrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func Encrypt(rsaPubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, nil)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func Decrypt(rsaPrivateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt data error: %w", err)
	}

	return decryptedData, nil
}
