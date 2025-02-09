package sha256mw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func SignRequestClientMiddleware(key string) ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		if next == nil {
			next = http.DefaultTransport
		}

		return requestSignerRoundTripper{next: next, key: key}
	}
}

type requestSignerRoundTripper struct {
	next http.RoundTripper
	key  string
}

func (rt requestSignerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// no key -> no signature
	if rt.key == "" {
		return rt.next.RoundTrip(r)
	}

	// read request body to calculate signature
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	_ = r.Body.Close()

	hash := hmac.New(sha256.New, []byte(rt.key))
	hash.Write(bodyBytes)
	hashSum := hash.Sum(nil)

	// already read request body -> gotta create new reader
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// pass signature in request headers
	r.Header.Set(sha256Header, base64.URLEncoding.EncodeToString(hashSum))

	return rt.next.RoundTrip(r)
}
