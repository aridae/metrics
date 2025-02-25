package rsacrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log"
)

// ExampleDecrypt демонстрирует базовую работу функции Decrypt.
func ExampleDecrypt() {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("ошибка генерации ключа: %v", err)
	}

	// Данные для шифрования
	data := []byte("привет, мир!")

	// Шифруем данные с использованием публичного ключа
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keyPair.PublicKey, data, nil)
	if err != nil {
		log.Fatalf("ошибка шифрования: %v", err)
	}

	// Расшифровка данных с использованием приватного ключа
	decryptedData, err := Decrypt(keyPair, encryptedData)
	if err != nil {
		log.Fatalf("ошибка расшифровки: %v", err)
	}

	// Вывод расшифрованных данных
	fmt.Printf("Расшифрованные данные: %s\n", decryptedData)

	// Output:
	// Расшифрованные данные: привет, мир!
}
