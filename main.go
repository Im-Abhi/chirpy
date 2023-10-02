package main

import (
	"log"
	"net/http"
)

func main() {
	const FILE_ROOT_PATH = "."
	const PORT = "8000"

	// create a new server mux
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(FILE_ROOT_PATH))))

	mux.HandleFunc("/healthz", handlerReadiness)	

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Handler: corsMux,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
