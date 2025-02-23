package rsamw

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	rsacrypto "github.com/aridae/go-metrics-store/pkg/rsa-crypto"
	"io"
	"net/http"
)

// EncryptRequestClientMiddleware создает middleware для HTTP-клиента, который шифрует тело запроса перед отправкой.
//
// Функция принимает публичный ключ RSA и возвращает функцию, которая оборачивает указанный RoundTripper.
// Если публичный ключ равен nil, middleware пропускает шифрование и передает запрос дальше без изменений.
//
// Параметры:
//
//	rsaPubKey *rsa.PublicKey — публичный ключ RSA для шифрования данных.
//
// Возвращаемое значение:
//
//	func(next http.RoundTripper) http.RoundTripper — middleware для HTTP-клиентов.
func EncryptRequestClientMiddleware(rsaPubKey *rsa.PublicKey) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		if next == nil {
			next = http.DefaultTransport
		}

		return encryptedRoundTripper{next: next, rsaPubKey: rsaPubKey}
	}
}

// encryptedRoundTripper реализует интерфейс http.RoundTripper, шифруя тело запроса перед передачей.
type encryptedRoundTripper struct {
	next      http.RoundTripper
	rsaPubKey *rsa.PublicKey
}

// RoundTrip реализует метод интерфейса http.RoundTripper, шифрующий тело запроса перед передачей.
func (rt encryptedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// no public key - no encryption
	if rt.rsaPubKey == nil {
		return rt.next.RoundTrip(r)
	}

	// read request body to calculate signature
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	_ = r.Body.Close()

	encryptedBodyBytes, err := rsacrypto.Encrypt(rt.rsaPubKey, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt body: %w", err)
	}

	// already read request body -> gotta create new reader
	r.Body = io.NopCloser(bytes.NewBuffer(encryptedBodyBytes))

	return rt.next.RoundTrip(r)
}
