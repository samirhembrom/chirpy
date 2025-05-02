package main

import (
	"context"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write(fmt.Appendf(nil, `Reset is only allowed in dev environment.`))
	}
	cfg.fileserverHits.Store(0)
	cfg.db.Reset(context.Background())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, "Hits reset to 0 and database reset to initial state."))
}
