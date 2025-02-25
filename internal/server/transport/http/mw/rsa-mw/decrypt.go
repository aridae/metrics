package rsamw

import (
	"bytes"
	"crypto/rsa"
	rsacrypto "github.com/aridae/go-metrics-store/pkg/rsa-crypto"
	"io"
	"net/http"
)

// DecryptRequestMiddleware создает middleware для HTTP-обработки, который расшифровывает зашифрованное тело запроса.
//
// Функция принимает приватный ключ RSA и возвращает обработчик, который выполняет следующие шаги:
// 1. Читает тело запроса.
// 2. Расшифровывает его с использованием предоставленного приватного ключа.
// 3. Заменяет тело запроса расшифрованными данными.
// 4. Продолжает выполнение следующего обработчика.
//
// Если предоставляется nil-указатель на приватный ключ, middleware пропускает расшифровку и сразу переходит к следующему обработчику.
//
// Параметры:
//
//	rsaPrivateKey *rsa.PrivateKey — приватный ключ RSA для расшифровки данных.
//
// Возвращаемое значение:
//
//	func(next http.Handler) http.Handler — middleware для обработки HTTP-запросов.
func DecryptRequestMiddleware(rsaPrivateKey *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rsaPrivateKey == nil {
				next.ServeHTTP(w, r)
				return
			}

			// read request body to calculate hashsum
			encryptedBodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()

			decryptedData, err := rsacrypto.Decrypt(rsaPrivateKey, encryptedBodyBytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// already read body -> gotta create new reader
			r.Body = io.NopCloser(bytes.NewBuffer(decryptedData))
			next.ServeHTTP(w, r)
		})
	}
}
