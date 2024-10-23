package handlers

import (
	"html/template"
	"net/http"

	metricstmplt "github.com/aridae/go-metrics-store/static/templates/metrics"
)

func (rt *Router) getAllMetricsHTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	metrics, err := rt.useCasesController.GetAllMetrics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pivotTableTmplt, err := template.ParseFS(metricstmplt.PivotTable, metricstmplt.PivotTableHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pivotTableTmplt.Execute(w, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
