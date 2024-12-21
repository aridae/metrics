package sha256mw

import "net/http"

const (
	sha256Header                          = "HashSHA256"
	requestForbiddenNoSignatureMessage    = "Forbidden: request must be signed"
	requestForbiddenWrongSignatureMessage = "Forbidden: wrong signature"
)

type ServerMiddleware func(next http.Handler) http.Handler

type ClientMiddleware func(next http.RoundTripper) http.RoundTripper
