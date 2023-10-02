package main

import (
	"log"
	"net/http"
	"fmt"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const FILE_ROOT_PATH = "."
	const PORT = "8000"

	// create instance of apiConfig
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// create a new server mux
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(FILE_ROOT_PATH)))))

	mux.HandleFunc("/healthz", handlerReadiness)	

	// handle the new handler
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	// handler reset hit count
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Handler: corsMux,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}