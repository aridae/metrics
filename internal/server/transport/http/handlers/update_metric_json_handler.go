package handlers

import (
	"encoding/json"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"net/http"
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

	metricFactory, err := resolveMetricFactoryForMetricType(transportMetric.MType)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	metric, err := buildMetricDomainModel(transportMetric, metricFactory)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}
	metricUpsertStrategy := metricFactory.ProvideUpsertStrategy()

	newMetricState, err := rt.useCasesController.UpsertScalarMetric(ctx, metric, metricUpsertStrategy)
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
