package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (rt *Router) getMetricByURLPathHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	metricTypeFromURL := chi.URLParam(r, urlParamMetricType)
	metricNameFromURL := chi.URLParam(r, urlParamMetricName)

	metricFactory, err := resolveMetricFactoryForMetricType(metricTypeFromURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricKey := metricFactory.CreateMetricKey(metricNameFromURL)

	latestState, err := rt.useCasesController.GetMetricByKey(ctx, metricKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if latestState == nil {
		http.Error(w, fmt.Sprintf("Metric %s not registered yet.", metricKey), http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(latestState.GetValue().String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
