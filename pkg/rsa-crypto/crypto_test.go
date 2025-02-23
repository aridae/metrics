package rsacrypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"math/big"
	"testing"
)

// Тестируем функцию Encrypt с валидными данными
func TestEncryptValid(t *testing.T) {
	// Создаем публичный ключ RSA
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsaPubKey := &rsaPrivateKey.PublicKey

	// Тестовые данные
	data := []byte("test data")

	// Вызываем функцию Encrypt
	encryptedData, err := Encrypt(rsaPubKey, data)

	// Проверяем отсутствие ошибок
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что зашифрованные данные не равны исходным данным
	if string(encryptedData) == string(data) {
		t.Error("encrypted data should not be equal to original data")
	}
}

// Тестируем функцию Encrypt с пустым массивом байтов
func TestEncryptEmptyData(t *testing.T) {
	// Создаем публичный ключ RSA
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsaPubKey := &rsaPrivateKey.PublicKey

	// Пустой массив байтов
	data := []byte{}

	// Вызываем функцию Encrypt
	_, err := Encrypt(rsaPubKey, data)

	// Ожидаемая ошибка
	expectedError := "empty data"

	// Проверяем наличие ошибки
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%v'", expectedError, err)
	}
}

// Тестируем функцию Encrypt с ошибкой шифрования
func TestEncryptEncryptionError(t *testing.T) {
	// Создание поддельного публичного ключа RSA с нулевым значением для вызова ошибки
	fakePublicKey := &rsa.PublicKey{N: big.NewInt(0)}

	// Тестовые данные
	data := []byte("test data")

	// Вызываем функцию Encrypt
	_, err := Encrypt(fakePublicKey, data)

	// Проверяем наличие ошибки
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// TestDecryptValid проверяет правильную работу функции Decrypt с валидными данными.
func TestDecryptValid(t *testing.T) {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Не удалось создать пару ключей RSA: %v", err)
	}

	// Данные для шифрования
	data := []byte("Тестовые данные")

	// Шифруем данные с использованием публичного ключа
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keyPair.PublicKey, data, nil)
	if err != nil {
		t.Fatalf("Ошибка шифрования: %v", err)
	}

	// Расшифровка данных с использованием приватного ключа
	decryptedData, err := Decrypt(keyPair, encryptedData)
	if err != nil {
		t.Fatalf("Ошибка расшифровки: %v", err)
	}

	// Проверяем, совпадают ли исходные и расшифрованные данные
	if !bytes.Equal(decryptedData, data) {
		t.Errorf("Расшифрованные данные не совпадают с исходными.\nОжидалось: %q\nПолучено: %q", data, decryptedData)
	}
}

// TestDecryptEmptyData проверяет поведение функции Decrypt при попытке расшифровать пустой массив байтов.
func TestDecryptEmptyData(t *testing.T) {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Не удалось создать пару ключей RSA: %v", err)
	}

	// Пробуем расшифровать пустой массив байтов
	_, err = Decrypt(keyPair, []byte{})

	// Проверяем, что функция возвращает ожидаемую ошибку
	if err == nil || err.Error() != "empty data" {
		t.Errorf("Ожидается ошибка 'empty data', получено: %v", err)
	}
}

// TestDecryptInvalidData проверяет поведение функции Decrypt при попытке расшифровать неправильно зашифрованные данные.
func TestDecryptInvalidData(t *testing.T) {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Не удалось создать пару ключей RSA: %v", err)
	}

	// Генерируем произвольные данные, которые не являются результатом шифрования OAEP
	badData := make([]byte, 1024)
	_, _ = rand.Read(badData)

	// Пробуем расшифровать эти данные
	_, err = Decrypt(keyPair, badData)

	// Проверяем, что функция возвращает ошибку
	if err == nil {
		t.Errorf("Ожидается ошибка, но она отсутствует")
	}
}
