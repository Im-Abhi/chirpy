package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
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

	// create a new router
	r := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	r.Get("/healthz", handlerReadiness)	

	// handle the new handler
	r.Get("/metrics", apiCfg.handlerMetrics)

	// handler reset hit count
	r.Get("/reset", apiCfg.handlerReset)

	corsRouter := middlewareCors(r)

	srv := &http.Server{
		Handler: corsRouter,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}
