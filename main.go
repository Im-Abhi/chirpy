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
	router := chi.NewRouter()


	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	// new api router
	apiRouter := chi.NewRouter()

	apiRouter.Get("/healthz", handlerReadiness)	
	// handle the new handler
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	// handler reset hit count
	apiRouter.Get("/reset", apiCfg.handlerReset)
	// mount the apiRouter router to r router through the /api route
	router.Mount("/api", apiRouter)

	corsRouter := middlewareCors(router)

	srv := &http.Server{
		Handler: corsRouter,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}
