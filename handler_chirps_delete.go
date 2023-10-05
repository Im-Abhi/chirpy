package main

import (
	"net/http"
	"strconv"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/Im-Abhi/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// when whether the user is authenticated or not
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	// if the user is authenticated get the chirpID to be deleted
	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp to be deleted")
		return
	}

	if err := CheckUser(userID, dbChirp.AuthorID); err != nil {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}

	// call the delete chirp function from the database
	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func CheckUser(userID, author_id int) (error) {
	if userID == author_id {
		return nil
	}

	return errors.New("Unauthorized User") 
}