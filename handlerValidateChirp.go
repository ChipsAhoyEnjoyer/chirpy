package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Println("Error with user request: Too many characters")
		respondWithError(w, http.StatusBadRequest, "Request too long; Max 140 char request.")
		return
	}
	payload := struct {
		CleanedBody string `json:"cleaned_body,omitempty"`
	}{CleanedBody: profanityFilter(params.Body)}
	respondWithJSON(w, http.StatusOK, payload)
}
