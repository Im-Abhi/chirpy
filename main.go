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
	// handler reset hit count
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Post("/validate_chirp", handlerChirpsValidate)
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

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	// struct for receiving json data
	type parameters struct {
		Body string `json:"body"`
	}

	// struct to return json value if valid
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// if there was some error decoding the json received
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	// if the chirp length is too long
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}