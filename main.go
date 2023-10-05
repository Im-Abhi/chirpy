package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/Im-Abhi/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret	   string
}

func main() {
	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()

	const FILE_ROOT_PATH = "."
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")

	if JWT_SECRET == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	db.ResetDB()

	// create instance of apiConfig
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB: db,
		jwtSecret: JWT_SECRET,
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

	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)

	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)	
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)	
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)	

	apiRouter.Post("/polka/webhooks", apiCfg.handlerWebhook)

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