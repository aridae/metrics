package handlers

import (
	"fmt"
	metricsfabrics "github.com/aridae/go-metrics-store/internal/server/metrics-fabrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (rt *Router) updateMetricByURLPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	if paramsCount := strings.Split(strings.Trim(r.URL.Path, "/"), "/"); len(paramsCount) != 4 {
		http.Error(w, "Unknown URL path.", http.StatusNotFound)
	}

	ctx := r.Context()

	metricTypeFromURL := chi.URLParam(r, urlParamMetricType)
	metricNameFromURL := chi.URLParam(r, urlParamMetricName)
	metricValueFromURL := chi.URLParam(r, urlParamMetricValue)

	metricFactory, err := resolveMetricFactoryFromURLPath(metricTypeFromURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricKey := metricFactory.CreateMetricKey(metricNameFromURL)
	metricValue, err := metricFactory.ParseScalarMetricValue(metricValueFromURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricToRegister := metricFactory.CreateScalarMetricToRegister(metricKey, metricValue)

	metricUpsertStrategy := metricFactory.ProvideUpsertStrategy()

	err = rt.useCasesController.UpsertScalarMetric(ctx, metricToRegister, metricUpsertStrategy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func resolveMetricFactoryFromURLPath(metricType string) (metricsfabrics.ScalarMetricFactory, error) {
	switch metricType {
	case counterURLParam:
		return metricsfabrics.NewCounterMetricFactory(), nil
	case gaugeURLParam:
		return metricsfabrics.NewGaugeMetricFactory(), nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}
}
