package handlers

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"net/http"
	"strconv"
	"strings"
)

const (
	indexMetricName         = 2
	indexMetricType         = 1
	indexMetricValue        = 3
	expectedPathParamsCount = 4

	counterURLParam = "counter"
	gaugeURLParam   = "gauge"
)

func getUpdateMetricByURLPathHandler(useCasesController useCasesController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
			return
		}

		params := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(params) < expectedPathParamsCount {
			http.Error(w, "unknown shit happened", http.StatusNotFound)
			return
		}

		metricUpdater, err := buildMetricUpdaterFromURLPath(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = useCasesController.UpsertScalarMetric(ctx, metricUpdater)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func buildMetricUpdaterFromURLPath(params []string) (models.ScalarMetricUpdater, error) {
	if len(params) < expectedPathParamsCount {
		return models.ScalarMetricUpdater{}, fmt.Errorf("invalid parameters count: expected at least %d, got %d", expectedPathParamsCount, len(params))
	}

	metricNameURLParam := params[indexMetricName]
	metricTypeURLParam := params[indexMetricType]
	metricValueURLParam := params[indexMetricValue]

	metricBuilderFn, ok := metricsConstructors[metricTypeURLParam]
	if !ok {
		return models.ScalarMetricUpdater{}, fmt.Errorf("unsupported metric type param: %s", metricTypeURLParam)
	}

	metric, err := metricBuilderFn(metricNameURLParam, metricValueURLParam)
	if err != nil {
		return models.ScalarMetricUpdater{}, fmt.Errorf("failed to build metric <type:%s> from provided params: %w", metricTypeURLParam, err)
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
