package sha256mw

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func SignResponseServerMiddleware(key string) ServerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(&responseSignerWriter{ResponseWriter: w, key: key}, r)
		})
	}
}

type responseSignerWriter struct {
	http.ResponseWriter
	key string
}

func (rw *responseSignerWriter) Write(bytes []byte) (int, error) {
	// no key -> no signature
	if rw.key == "" {
		return rw.ResponseWriter.Write(bytes)
	}

	hash := hmac.New(sha256.New, []byte(rw.key))
	hash.Write(bytes)
	hashSum := hash.Sum(nil)

	rw.Header().Set(sha256Header, base64.URLEncoding.EncodeToString(hashSum))

	return rw.ResponseWriter.Write(bytes)
}
