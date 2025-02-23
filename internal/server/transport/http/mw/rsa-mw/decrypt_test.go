package rsamw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDecryptRequestMiddlewareNilKey проверяет, что middleware пропускает расшифровку, если ключ равен nil.
func TestDecryptRequestMiddlewareNilKey(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("test body"))
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := DecryptRequestMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write(body)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}

	respBody, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(respBody), "test body") {
		t.Errorf("response does not contain test body")
	}
}

// TestDecryptRequestMiddlewareSuccess проверяет успешную расшифровку тела запроса.
func TestDecryptRequestMiddlewareSuccess(t *testing.T) {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Не удалось создать пару ключей RSA: %v", err)
	}

	plainText := []byte("hello world")
	cipherText, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keyPair.PublicKey, plainText, nil)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(cipherText))
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := DecryptRequestMiddleware(keyPair)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write(body)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}

	respBody, _ := io.ReadAll(rr.Body)
	if !bytes.Equal(respBody, plainText) {
		t.Errorf("expected response body to match plaintext, got %s", respBody)
	}
}

// TestDecryptRequestMiddlewareError проверяет поведение middleware при возникновении ошибки расшифровки.
func TestDecryptRequestMiddlewareError(t *testing.T) {
	// Генерируем пару ключей RSA
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Не удалось создать пару ключей RSA: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("invalid cipher text"))
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := DecryptRequestMiddleware(keyPair)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status Internal Server Error, got %v", rr.Code)
	}
}
