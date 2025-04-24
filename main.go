package main

import (
	"log"
	"net/http"
)

func main() {
	const fileRootPath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(fileRootPath)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s", fileRootPath, port)
	log.Fatal(srv.ListenAndServe())
}
