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
		respondWithError(w, http.StatusInternalServerError, "Error decoding response: "+err.Error())
		return
	}
	u, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error posting email '%v' to database: %v\n", req.Email, err)
		respondWithError(w, http.StatusInternalServerError, "Error registering user: "+err.Error())
		return
	}
	resp := user{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}
