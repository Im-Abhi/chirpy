package main

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	
	user, err := cfg.DB.CreateUser(params.Email, string(hashedPassword))
	// if error is there -> could be possible that email already exists
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		ID int `json:"id"`
		Email string `json:"email"`
	} {
		ID:   user.ID,
		Email: user.Email,
	})
}