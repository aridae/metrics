package handlers

import (
	"encoding/json"
	"net/http"

	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func (rt *Router) updateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	transportMetric := httpmodels.Metric{}
	err := json.NewDecoder(r.Body).Decode(&transportMetric)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	metric, err := buildMetricDomainModel(transportMetric)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	newMetricState, err := rt.useCasesController.UpsertMetric(ctx, metric)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	transportMetricAfterUpsert, err := buildMetricTransportModel(newMetricState)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transportMetricAfterUpsert)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}
}
