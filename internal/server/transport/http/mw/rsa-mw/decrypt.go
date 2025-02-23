package rsamw

import (
	"bytes"
	"crypto/rsa"
	rsacrypto "github.com/aridae/go-metrics-store/pkg/rsa-crypto"
	"io"
	"net/http"
)

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
