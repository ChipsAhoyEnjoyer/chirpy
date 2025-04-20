package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Function currently assumes user request is valid without checking
	defer r.Body.Close()
	req := struct {
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding response")
	}
	resp, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error adding email '%v' to database: %v\n", req.Email, err)
		respondWithError(w, http.StatusInternalServerError, "Error registering user: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, resp)
}
