package sha256mw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

/*
responseSignerWriter{ResponseWriter: w, key: key}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()

			hash := hmac.New(sha256.New, []byte(key))
			hash.Write(bodyBytes)
			hashSum := hash.Sum(nil)

			// already read body -> gotta create new reader
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			r.Header.Set(Header, string(hashSum))
*/

func ValidateRequestServerMiddleware(key string) ServerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// no key -> no signature
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			// get signature from headers
			signature := r.Header.Get(sha256Header)
			if signature == "" {
				http.Error(w, requestForbiddenNoSignatureMessage, http.StatusForbidden)
				return
			}

			signatureBytes, err := base64.URLEncoding.DecodeString(signature)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// read request body to calculate hashsum
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()

			hash := hmac.New(sha256.New, []byte(key))
			hash.Write(bodyBytes)
			hashSum := hash.Sum(nil)

			// check: requst hashsum must be equal to header signature
			if !hmac.Equal(hashSum, signatureBytes) {
				http.Error(w, requestForbiddenWrongSignatureMessage, http.StatusBadRequest)
				return
			}

			// already read body -> gotta create new reader
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
		})
	}
}
