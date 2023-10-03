package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/Im-Abhi/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const FILE_ROOT_PATH = "."
	const PORT = "8000"

	db, _ := database.NewDB("database.json")
	// create instance of apiConfig
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB: db,
	}

	// create a new router
	router := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	// new api router
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)	
	// handler reset hit count
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)	
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)	

	apiRouter.Post("/users", apiCfg.handlerUsersCreate)

	// mount the apiRouter router to r router through the /api route
	router.Mount("/api", apiRouter)
	// metric router
	adminRouter := chi.NewRouter()

	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	router.Mount("/admin", adminRouter)

	corsRouter := middlewareCors(router)

	srv := &http.Server{
		Handler: corsRouter,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}