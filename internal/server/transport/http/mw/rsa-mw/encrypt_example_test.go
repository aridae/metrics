package rsamw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ExampleEncryptRequestClientMiddleware показывает, как работает middleware для шифрования тела запроса перед отправкой.
func ExampleEncryptRequestClientMiddleware() {
	// Генерация публичного ключа RSA
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("ошибка генерации ключа: %v", err)
	}

	// Создание HTTP-запроса с открытым текстом
	plainText := []byte("hello world")
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(plainText))
	if err != nil {
		log.Fatalf("ошибка создания запроса: %v", err)
	}

	// Настройка middleware и выполнение запроса
	rt := EncryptRequestClientMiddleware(&key.PublicKey)(roundTripStub{})
	resp, err := rt.RoundTrip(req)
	if err != nil {
		log.Fatalf("ошибка отправки запроса: %v", err)
	}
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	// Получение зашифрованного тела запроса
	body, _ := io.ReadAll(req.Body)
	fmt.Println("Зашифрованное тело запроса:", string(body))
}
