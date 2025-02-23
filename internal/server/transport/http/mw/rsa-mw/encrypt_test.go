package rsamw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type roundTripStub struct{}

func (s roundTripStub) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

// TestEncryptRequestClientMiddlewareNilKey проверяет, что middleware пропускает шифрование, если ключ равен nil.
func TestEncryptRequestClientMiddlewareNilKey(t *testing.T) {
	plainText := []byte("hello world")

	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(plainText))

	rt := EncryptRequestClientMiddleware(nil)(roundTripStub{})

	_, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что тело запроса НЕ было зашифровано
	body, _ := io.ReadAll(req.Body)
	require.Equal(t, plainText, body)
}

// TestEncryptRequestClientMiddlewareSuccess проверяет успешное шифрование тела запроса.
func TestEncryptRequestClientMiddlewareSuccess(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	plainText := []byte("hello world")
	cipherText, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, plainText, nil)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(plainText))

	rt := EncryptRequestClientMiddleware(&key.PublicKey)(roundTripStub{})

	_, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что тело запроса было зашифровано
	body, _ := io.ReadAll(req.Body)
	if reflect.DeepEqual(body, cipherText) {
		t.Errorf("body was not encrypted")
	}
}

// TestEncryptRequestClientMiddlewareError проверяет поведение middleware при возникновении ошибки шифрования.
func TestEncryptRequestClientMiddlewareError(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Создаем невалидный запрос
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}

	rt := EncryptRequestClientMiddleware(&key.PublicKey)(nil)

	_, err = rt.RoundTrip(req)
	require.Error(t, err)
}
