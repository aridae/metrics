package handlers

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (rt *Router) getMetricByURLPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	metricName := chi.URLParam(r, urlParamMetricName)
	metricType, err := mapAPIToDomainScalarMetricType(chi.URLParam(r, urlParamMetricType))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricKey := models.BuildMetricKey(metricType, metricName)

	latestState, err := rt.useCasesController.GetScalarMetricLatestState(ctx, metricKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if latestState == nil {
		http.Error(w, fmt.Sprintf("Metric %s:%s not registered yet.", metricType, metricName), http.StatusNotFound)
		return
	}
	metricValueStr := fmt.Sprintf("%v", latestState.Value)

	_, err = w.Write([]byte(metricValueStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func mapAPIToDomainScalarMetricType(apiT string) (models.ScalarMetricType, error) {
	switch apiT {
	case counterURLParam:
		return models.ScalarMetricTypeCounter, nil
	case gaugeURLParam:
		return models.ScalarMetricTypeGauge, nil
	default:
		return "", fmt.Errorf("unknown scalar metrics type: %s", apiT)
	}
}
