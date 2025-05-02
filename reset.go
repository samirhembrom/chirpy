package main

import (
	"context"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(
			w,
			http.StatusForbidden,
			"You don't have access",
			errors.New("Non dev user tring to access"),
		)
	}
	err := cfg.db.DeleteUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
}
