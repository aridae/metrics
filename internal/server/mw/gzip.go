package mw

import (
	"compress/gzip"
	"github.com/aridae/go-metrics-store/internal/server/logger"
	"net/http"
	"strings"
)

func GzipDecompressRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		uncompressedBody := r.Body
		gzipReader, err := gzip.NewReader(uncompressedBody)
		if err != nil {
			logger.Obtain().Errorf("[mw.GzipDecompressRequestMiddleware] failed to create gzip.Reader: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uncompressedBody.Close()

		r.Body = gzipReader
		next.ServeHTTP(w, r)
	})
}

func GzipCompressResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// is like gzip.NewWriter but specifies the compression level instead
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, compressor: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	compressor *gzip.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.compressor.Write(b)
}
