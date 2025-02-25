package rsacrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Encrypt шифрует данные с использованием алгоритма RSA-OAEP.
//
// Функция принимает публичный ключ RSA и данные для шифрования.
// Возвращается зашифрованные данные и возможная ошибка.
//
// Если входные данные пусты, возвращается ошибка "empty data".
// Если возникает ошибка при шифровании, возвращается соответствующая ошибка.
//
// Пример использования:
//
//	pk, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//		log.Fatalf("failed to generate key: %v", err)
//	}
//	encryptedData, err := Encrypt(pk, []byte("secret message"))
//	if err != nil {
//		log.Fatalf("failed to encrypt: %v", err)
//	}
//	fmt.Printf("Encrypted Data: %x\n", encryptedData)
//
// Параметры:
//
//	rspPubKey *rsa.PublicKey — публичный ключ RSA для шифрования.
//	data []byte — данные для шифрования.
//
// Возвращаемые значения:
//
//	[]byte — зашифрованные данные.
//	error — ошибка, если произошла ошибка при шифровании.
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

// Decrypt расшифровывает данные, зашифрованные с использованием алгоритма RSA-OAEP.
//
// Функция принимает приватный ключ RSA и зашифрованные данные.
// Возвращает расшифрованные данные и возможную ошибку.
//
// Если входные данные пусты, возвращается ошибка "empty data".
// Если возникает ошибка при расшифровке, возвращается соответствующая ошибка.
//
// Пример использования:
//
//	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//		log.Fatalf("failed to generate key: %v", err)
//	}
//	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keyPair.PublicKey, []byte("secret message"), nil)
//	if err != nil {
//		log.Fatalf("failed to encrypt: %v", err)
//	}
//	decryptedData, err := Decrypt(&keyPair, encryptedData)
//	if err != nil {
//		log.Fatalf("failed to decrypt: %v", err)
//	}
//	fmt.Printf("Decrypted Data: %s\n", decryptedData)
//
// Параметры:
//
//	rsaPrivateKey *rsa.PrivateKey — приватный ключ RSA для расшифровки.
//	encryptedData []byte — зашифрованные данные.
//
// Возвращаемые значения:
//
//	[]byte — расшифрованные данные.
//	error — ошибка, если произошла ошибка при расшифровке.
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
