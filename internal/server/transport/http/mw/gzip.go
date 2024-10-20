package mw

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/aridae/go-metrics-store/internal/server/logger"
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

		next.ServeHTTP(&gzipWriter{ResponseWriter: w, compressor: gz}, r)
	})
}

type gzipWriter struct {
	wroteHeader bool

	http.ResponseWriter
	compressor io.Writer
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	if !gw.wroteHeader {
		gw.WriteHeader(http.StatusOK)
	}

	return gw.compressor.Write(b)
}

func (gw *gzipWriter) WriteHeader(code int) {
	if gw.wroteHeader {
		return
	}
	gw.wroteHeader = true
	defer gw.ResponseWriter.WriteHeader(code)

	gw.Header().Set("Content-Encoding", "gzip")
	gw.Header().Del("Content-Length")
}
