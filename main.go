package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8000"

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /healthz", handlerRediness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerResetMetrics)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	log.Println("Server started on port", port)
	log.Fatal(server.ListenAndServe())
}
