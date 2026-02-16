package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in dev environment", nil)
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error(), err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits reset to %d and database reset to initial state", cfg.fileserverHits.Load())
}
