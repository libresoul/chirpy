package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits were reset"))
}
