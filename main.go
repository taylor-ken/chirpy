package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	fileserverHits := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", fileserverHits.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", fileserverHits.handlerMetrics)
	mux.HandleFunc("/reset", fileserverHits.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	response := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(response))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	response := fmt.Sprintf("Hits: %v", count)
	w.Write([]byte(response))
}
