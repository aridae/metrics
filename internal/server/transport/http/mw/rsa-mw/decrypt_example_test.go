package rsamw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// ExampleDecryptRequestMiddleware показывает, как работает middleware для расшифровки зашифрованного тела запроса.
func ExampleDecryptRequestMiddleware() {
	// Генерация приватного ключа RSA
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("ошибка генерации ключа: %v", err)
	}

	// Шифрование данных
	plainText := []byte("hello world")
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, plainText, nil)
	if err != nil {
		log.Fatalf("ошибка шифрования: %v", err)
	}

	// Создание HTTP-запроса с зашифрованным телом
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(cipherText))
	if err != nil {
		log.Fatalf("ошибка создания запроса: %v", err)
	}

	// Настройка middleware и обработчика
	rr := httptest.NewRecorder()
	handler := DecryptRequestMiddleware(key)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, "Расшифрованное тело: %s", body)
	}))

	// Выполнение запроса
	handler.ServeHTTP(rr, req)

	// Проверка результата
	respBody, _ := io.ReadAll(rr.Body)
	fmt.Println("Ответ от сервера:", string(respBody))

	// Output:
	// Ответ от сервера: Расшифрованное тело: hello world
}
