package mw

import (
	"github.com/aridae/go-metrics-store/internal/server/logger"
	"net/http"
	"strings"
)

func TrimTrailingSlash(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			logger.Obtain().Debugf("trimming trailing slash in uri %s", r.URL.Path)

			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
