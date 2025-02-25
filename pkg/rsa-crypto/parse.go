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

// ParsePublicKey парсит PEM-кодированный публичный ключ RSA.
//
// Функция принимает массив байтов, представляющий PEM-кодированный публичный ключ,
// и пытается распарсить его в структуру rsa.PublicKey.
//
// Если входные данные не содержат корректного PEM-блока или если не удается распознать
// публичный ключ как ключ RSA, возвращается ошибка.
//
// Пример использования:
//
//	publicKeyPEM := []byte(`-----BEGIN PUBLIC KEY-----
//	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCfT9kBoLavP0bKDOOZxPjJFOaU
//	3E1hOvm15ATCG6Q+nFjbBz2GG2pOKPLY3TVC2UW9XfthdS5oD8UDLzVOr37RrbkR
//	FIl1y5WcZZPPiCJmlE3XRe6g9WK2T6xtfCeOS3gkv1/HY6GLG/bhBzrV1Q2k7pbVU
//	G+vSHFYZEnSQIDAQAB
//	-----END PUBLIC KEY-----`)
//	pubKey, err := ParsePublicKey(publicKeyPEM)
//	if err != nil {
//		log.Fatalf("не удалось распарсить публичный ключ: %v", err)
//	}
//	fmt.Println("успешно распарсирован публичный ключ:", pubKey.N.String())
//
// Параметры:
//
//	publicKeyData []byte — массив байтов, содержащий PEM-кодированный публичный ключ.
//
// Возвращаемые значения:
//
//	*rsa.PublicKey — успешно распарсенный публичный ключ RSA.
//	error — ошибка, если произошла ошибка при парсинге.
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
		return nil, fmt.Errorf("unsupported type of public key: %T", pubKey)
	}

	return rsaPubKey, nil
}

// ParsePrivateKey парсит PEM-кодированный приватный ключ RSA.
//
// Функция принимает массив байтов, представляющий PEM-кодированный приватный ключ,
// и пытается распарсить его в структуру rsa.PrivateKey.
//
// Если входные данные не содержат корректного PEM-блока или если не удается распознать
// приватный ключ как ключ RSA, возвращается ошибка.
//
// Пример использования:
//
//	privateKeyPEM := []byte(`-----BEGIN PRIVATE KEY-----
//	MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgSRL6DsfzYAxuEg7N
//	jUdu3ikO4wXjltWl/Nh8HCwe1JhRANCAARLHtrZVdJ9q4AKshwKnbsxBVKc/rUBp
//	LJbVMCRCh3InqTMy2IJvxIY5lgWbbi27DDkPl1cKZr3udL4eqBZRRW
//	-----END PRIVATE KEY-----`)
//	privateKey, err := ParsePrivateKey(privateKeyPEM)
//	if err != nil {
//		log.Fatalf("не удалось распарсить приватный ключ: %v", err)
//	}
//	fmt.Println("успешно распарсирован приватный ключ:", privateKey.D.Int.Text(16))
//
// Параметры:
//
//	privateKeyData []byte — массив байтов, содержащий PEM-кодированный приватный ключ.
//
// Возвращаемые значения:
//
//	*rsa.PrivateKey — успешно распарсированный приватный ключ RSA.
//	error — ошибка, если произошла ошибка при парсинге.
func ParsePrivateKey(privateKeyData []byte) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode(privateKeyData)
	if pemBlock == nil {
		return nil, errors.New("no PEM block found")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unsupported type of private key: %T", privateKey)
	}

	return rsaPrivateKey, nil
}

// FromFile загружает RSA-ключ (приватный или публичный) из файла.
//
// Функция принимает путь к файлу и функцию для парсинга байтового представления ключа.
// Если файл не найден или не удается прочитать его содержимое, возвращается ошибка.
// Если функция парсинга завершится неудачей, также возвращается ошибка.
//
// Пример использования:
//
//	pathToPublicKey := "/path/to/public_key.pem"
//	publicKey, err := FromFile[pathToPublicKey, ParsePublicKey]
//	if err != nil {
//		log.Fatalf("failed to load public key from file: %v", err)
//	}
//	fmt.Println("successfully loaded public key:", publicKey.N.String())
//
// pathToPrivateKey := "/path/to/private_key.pem"
//
//	privateKey, err := FromFile[pathToPrivateKey, ParsePrivateKey]
//	if err != nil {
//		log.Fatalf("failed to load private key from file: %v", err)
//	}
//	fmt.Println("successfully loaded private key:", privateKey.D.Int.Text(16))
//
// Параметры:
//
//	path string — путь к файлу, содержащему ключ.
//	parseFn func([]byte) (*Key, error) — функция для парсинга байтового представления ключа.
//
// Возвращаемые значения:
//
//	*Key — успешно загруженный ключ.
//	error — ошибка, если возникла проблема при загрузке или парсинге ключа.
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
