package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"net/http"
)

func (rt *Router) updateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	jsonMetric := httpmodels.MetricUpsert{}

	err := json.NewDecoder(r.Body).Decode(&jsonMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonMetricValue, err := resolveMetricValueFromJSONMetric(jsonMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricFactory, err := resolveMetricFactoryForMetricType(jsonMetric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricKey := metricFactory.CreateMetricKey(jsonMetric.ID)
	metricValue, err := metricFactory.CastScalarMetricValue(jsonMetricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric := metricFactory.CreateScalarMetricToRegister(metricKey, metricValue)
	metricUpsertStrategy := metricFactory.ProvideUpsertStrategy()

	err = rt.useCasesController.UpsertScalarMetric(ctx, metric, metricUpsertStrategy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func resolveMetricValueFromJSONMetric(metric httpmodels.MetricUpsert) (any, error) {
	switch metric.MType {
	case gauge:
		if metric.Value == nil {
			return nil, errors.New("'value' field is required for gauge metric")
		}
		return *metric.Value, nil
	case counter:
		if metric.Delta == nil {
			return nil, errors.New("'delta' field is required for counter metric")
		}
		return *metric.Delta, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
