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
	// handle / route to serve a static html file from the root directory
	mux.Handle("/", http.FileServer(http.Dir(FILE_ROOT_PATH)))
	// serving the logo url path should match with the directory path
	mux.Handle("/assets/", http.FileServer(http.Dir(FILE_ROOT_PATH)))
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Handler: corsMux,
		Addr: 	":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	log.Printf("Serving on port : %s\n", PORT)
	// listen and serve
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
