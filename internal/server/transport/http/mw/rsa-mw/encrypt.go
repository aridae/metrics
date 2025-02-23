package rsamw

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	rsacrypto "github.com/aridae/go-metrics-store/pkg/rsa-crypto"
	"io"
	"net/http"
)

func EncryptRequestClientMiddleware(rsaPubKey *rsa.PublicKey) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		if next == nil {
			next = http.DefaultTransport
		}

		return encryptedRoundTripper{next: next, rsaPubKey: rsaPubKey}
	}
}

type encryptedRoundTripper struct {
	next      http.RoundTripper
	rsaPubKey *rsa.PublicKey
}

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
