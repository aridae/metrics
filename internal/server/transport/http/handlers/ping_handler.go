package handlers

import "net/http"

func (rt *Router) pingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	for _, dep := range rt.checkIfAvailableOnPing {
		if err := dep.Ping(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
