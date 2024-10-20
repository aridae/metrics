package handlers

import (
	"encoding/json"
	"net/http"

	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func (rt *Router) getMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	transportMetricRequest := httpmodels.MetricRequest{}
	err := json.NewDecoder(r.Body).Decode(&transportMetricRequest)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	metricFactory, err := resolveMetricFactoryForMetricType(transportMetricRequest.MType)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}
	metricKey := metricFactory.CreateMetricKey(transportMetricRequest.ID)

	metric, err := rt.useCasesController.GetMetricByKey(ctx, metricKey)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	if metric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	transportMetric, err := buildMetricTransportModel(*metric)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transportMetric)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}
}

func mustWriteJSONError(w http.ResponseWriter, err error, code int) {
	errMsg := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}

	errMessagePayload, _ := json.Marshal(errMsg)

	_, _ = w.Write(errMessagePayload)
	w.WriteHeader(code)
}
