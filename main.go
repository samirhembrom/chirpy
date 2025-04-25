package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const fileRootPath = "."

const port = "8080"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("Request received, hits now: %d", cfg.fileserverHits.Load())
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileRootPath)))),
	)
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", cfg.handlerWriteRequests())
	mux.HandleFunc("/reset", cfg.handlerResetRequests())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s", fileRootPath, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerWriteRequests() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var count int32 = cfg.fileserverHits.Load()
		str := fmt.Sprintf("Hits: %d", count)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(str))
	}
}

func (cfg *apiConfig) handlerResetRequests() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Store(0)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
	}
}
