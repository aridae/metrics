package loggingmw

import (
	"github.com/aridae/go-metrics-store/pkg/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		respTracker := &responseTracker{ResponseWriter: w}

		next.ServeHTTP(respTracker, r)

		duration := time.Since(start)

		logger.Infof("[mw.LoggingMiddleware] call %s %s took %s, handled with status code %d, resp body size %d bytes",
			r.Method, r.RequestURI, duration, respTracker.status, respTracker.sizeBytes,
		)
	})
}

type responseTracker struct {
	http.ResponseWriter
	status      int
	sizeBytes   int
	wroteHeader bool
}

func (o *responseTracker) Write(bytes []byte) (int, error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}

	size, err := o.ResponseWriter.Write(bytes)
	o.sizeBytes = size

	return size, err
}

func (o *responseTracker) WriteHeader(code int) {
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	defer o.ResponseWriter.WriteHeader(code)

	o.status = code
}
