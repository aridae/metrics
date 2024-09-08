package handlers

import (
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases"
	"net/http"
)

const (
	urlPathUpdate = "/update/"
)

func NewRouter(useCasesController *usecases.Controller) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(urlPathUpdate, getUpdateMetricByURLPathHandler(useCasesController))

	return mux
}
