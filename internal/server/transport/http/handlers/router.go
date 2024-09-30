package handlers

import (
	"context"
	"fmt"
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	urlParamMetricType  = "metric_type"
	urlParamMetricName  = "metric_name"
	urlParamMetricValue = "metric_value"

	counterURLParam = "counter"
	gaugeURLParam   = "gauge"
)

var (
	updateMetricValueURLPath  = fmt.Sprintf("/update/{%s}/{%s}/{%s}", urlParamMetricType, urlParamMetricName, urlParamMetricValue)
	getMetricValueURLPath     = fmt.Sprintf("/value/{%s}/{%s}", urlParamMetricType, urlParamMetricName)
	getAllMetricValuesURLPath = "/"
)

type useCasesController interface {
	UpsertScalarMetric(ctx context.Context, metricToRegister models.ScalarMetricToRegister, strategy metricsupsertstrategies.Strategy) error
	GetScalarMetricLatestState(ctx context.Context, metricKey models.MetricKey) (*models.ScalarMetric, error)
	GetAllScalarMetricsLatestStates(ctx context.Context) ([]models.ScalarMetric, error)
}

type Router struct {
	useCasesController useCasesController
	httpMux            *chi.Mux
}

func NewRouter(useCasesController useCasesController) *Router {
	chiMux := chi.NewRouter()

	router := &Router{
		useCasesController: useCasesController,
		httpMux:            chiMux,
	}

	chiMux.HandleFunc(updateMetricValueURLPath, router.updateMetricByURLPathHandler)
	chiMux.HandleFunc(getMetricValueURLPath, router.getMetricByURLPathHandler)
	chiMux.HandleFunc(getAllMetricValuesURLPath, router.getAllMetricsHTMLHandler)

	return router
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.httpMux.ServeHTTP(w, r)
}
