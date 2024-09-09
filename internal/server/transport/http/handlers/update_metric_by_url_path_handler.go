package handlers

import (
	"errors"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (rt *Router) updateMetricByURLPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	metricType := chi.URLParam(r, urlParamMetricType)
	metricName := chi.URLParam(r, urlParamMetricName)
	metricValue := chi.URLParam(r, urlParamMetricValue)

	metricUpdater, err := buildMetricUpdaterFromURLPath(metricType, metricName, metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.useCasesController.UpsertScalarMetric(ctx, metricUpdater)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func buildMetricUpdaterFromURLPath(metricType string, metricName string, metricValue string) (models.ScalarMetricUpdater, error) {
	if metricType == "" || metricName == "" || metricValue == "" {
		return models.ScalarMetricUpdater{}, errors.New("invalid parameters: got empty string/s")
	}

	metricBuilderFn, ok := metricsConstructors[metricType]
	if !ok {
		return models.ScalarMetricUpdater{}, fmt.Errorf("unsupported scalar-metrics type param: %s", metricType)
	}

	metric, err := metricBuilderFn(metricName, metricValue)
	if err != nil {
		return models.ScalarMetricUpdater{}, fmt.Errorf("failed to build scalar-metrics <type:%s> from provided params: %w", metricType, err)
	}

	return metric, nil
}

var metricsConstructors = map[string]func(name, value string) (models.ScalarMetricUpdater, error){
	counterURLParam: buildCounter,
	gaugeURLParam:   buildGauge,
}

func buildCounter(name, value string) (models.ScalarMetricUpdater, error) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return models.ScalarMetricUpdater{}, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return models.ScalarMetricUpdater{
		Type:  models.ScalarMetricTypeCounter,
		Name:  name,
		Value: parsedValue,
	}, nil
}

func buildGauge(name, value string) (models.ScalarMetricUpdater, error) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return models.ScalarMetricUpdater{}, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	return models.ScalarMetricUpdater{
		Type:  models.ScalarMetricTypeGauge,
		Name:  name,
		Value: parsedValue,
	}, nil
}
