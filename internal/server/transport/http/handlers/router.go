package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/go-chi/chi/v5"
)

const (
	urlParamMetricType  = "metric_type"
	urlParamMetricName  = "metric_name"
	urlParamMetricValue = "metric_value"

	counter = "counter"
	gauge   = "gauge"
)

var (
	updateMetricWithJSONBodyURLPath       = "/update"
	updateMetricsBatchWithJSONBodyURLPath = "/updates/"
	getMetricWithJSONBodyURLPath          = "/value"
	updateMetricWithURLParamsValueURLPath = fmt.Sprintf("/update/{%s}/{%s}/{%s}", urlParamMetricType, urlParamMetricName, urlParamMetricValue)
	getMetricValueURLPath                 = fmt.Sprintf("/value/{%s}/{%s}", urlParamMetricType, urlParamMetricName)
	getAllMetricValuesURLPath             = "/"
	pingHandlerURLPath                    = "/ping"
)

type pingable interface {
	Ping(context.Context) error
}

type useCasesController interface {
	UpsertMetricsBatch(context.Context, []models.MetricUpsert) ([]models.Metric, error)
	UpsertMetric(context.Context, models.MetricUpsert) (models.Metric, error)
	GetMetricByKey(context.Context, models.MetricKey) (*models.Metric, error)
	GetAllMetrics(context.Context) ([]models.Metric, error)
}

type Router struct {
	useCasesController useCasesController
	httpMux            *chi.Mux

	checkIfAvailableOnPing []pingable
}

func NewRouter(useCasesController useCasesController, options ...RouterOption) *Router {
	opts := new(routerOpts)
	for _, applyOption := range options {
		applyOption(opts)
	}

	chiMux := chi.NewRouter()

	router := &Router{
		useCasesController:     useCasesController,
		httpMux:                chiMux,
		checkIfAvailableOnPing: opts.checkAvailableOnPing,
	}

	chiMux.HandleFunc(updateMetricWithJSONBodyURLPath, router.updateMetricJSONHandler)
	chiMux.HandleFunc(updateMetricWithJSONBodyURLPath+"/", router.updateMetricJSONHandler) // trailing slash

	chiMux.HandleFunc(updateMetricsBatchWithJSONBodyURLPath, router.updateMetricsBatchJSONHandler)
	chiMux.HandleFunc(updateMetricsBatchWithJSONBodyURLPath+"/", router.updateMetricsBatchJSONHandler) // trailing slash

	chiMux.HandleFunc(updateMetricWithURLParamsValueURLPath, router.updateMetricByURLPathHandler)

	chiMux.HandleFunc(getMetricWithJSONBodyURLPath, router.getMetricJSONHandler)
	chiMux.HandleFunc(getMetricWithJSONBodyURLPath+"/", router.getMetricJSONHandler) // trailing slash

	chiMux.HandleFunc(pingHandlerURLPath, router.pingHandler)
	chiMux.HandleFunc(pingHandlerURLPath+"/", router.pingHandler) // trailing slash

	chiMux.HandleFunc(getMetricValueURLPath, router.getMetricByURLPathHandler)

	chiMux.HandleFunc(getAllMetricValuesURLPath, router.getAllMetricsHTMLHandler)

	if opts.serveDebugPprof {
		chiMux.Mount(opts.debugPprofPattern, http.DefaultServeMux)
	}

	return router
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.httpMux.ServeHTTP(w, r)
}

type routerOpts struct {
	checkAvailableOnPing []pingable
	serveDebugPprof      bool
	debugPprofPattern    string
}

type RouterOption func(opts *routerOpts)

func CheckAvailableOnPing(dep pingable) RouterOption {
	return func(opts *routerOpts) {
		opts.checkAvailableOnPing = append(opts.checkAvailableOnPing, dep)
	}
}

func WithDebugPprof(pattern string) RouterOption {
	return func(opts *routerOpts) {
		opts.serveDebugPprof = true
		opts.debugPprofPattern = pattern
	}
}
