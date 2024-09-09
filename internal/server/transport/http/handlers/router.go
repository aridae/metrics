package handlers

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"net/http"
)

const (
	urlPathUpdate = "/update/"
)

type useCasesController interface {
	UpsertScalarMetric(ctx context.Context, updater models.ScalarMetricUpdater) error
}

func NewRouter(useCasesController useCasesController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(urlPathUpdate, getUpdateMetricByURLPathHandler(useCasesController))

	return mux
}
