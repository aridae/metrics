package handlers

import (
	"encoding/json"
	"net/http"

	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"github.com/aridae/go-metrics-store/pkg/slice"
)

func (rt *Router) updateMetricsBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	var transportMetrics httpmodels.Metrics
	err := json.NewDecoder(r.Body).Decode(&transportMetrics)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	err = transportMetrics.Validate()
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	metricsUpserts, err := slice.MapBatch(transportMetrics, buildMetricDomainModel)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	upsertedMetrics, err := rt.useCasesController.UpsertMetricsBatch(ctx, metricsUpserts)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	transportUpsertedMetrics, err := slice.MapBatch(upsertedMetrics, buildMetricTransportModel)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transportUpsertedMetrics)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}
}
