package handlers

import (
	"encoding/json"
	"fmt"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"net/http"
)

func (rt *Router) getMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	transportMetricRequest := httpmodels.MetricRequest{}
	err := json.NewDecoder(r.Body).Decode(&transportMetricRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricFactory, err := resolveMetricFactoryForMetricType(transportMetricRequest.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricKey := metricFactory.CreateMetricKey(transportMetricRequest.ID)

	metric, err := rt.useCasesController.GetScalarMetricLatestState(ctx, metricKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metric == nil {
		http.Error(w, fmt.Sprintf("Metric %s not registered yet.", metricKey), http.StatusNotFound)
		return
	}

	transportMetric, err := buildMetricTransportModel(*metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transportMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
