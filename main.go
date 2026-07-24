package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func _(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileServerHits atomic.Int32
}

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

func main() {
	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", handlerRediness)
	mux.HandleFunc("/metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset/", apiCfg.handlerResetMetrics)

	server := &http.Server{Addr: ":8080", Handler: mux}
	log.Println("Server started on port 8080")
	log.Fatal(server.ListenAndServe())
}

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatal("Failed to write response: ", err)
	}
}
