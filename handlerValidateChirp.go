package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding response")
		return
	}

	if len(req.Body) > 140 {
		log.Println("Error with user request: Too many characters")
		respondWithError(w, http.StatusBadRequest, "Request too long; Max 140 char request.")
		return
	}

	cleanedChirp := profanityFilter(req.Body)
	c, err := cfg.dbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedChirp,
			UserID: req.UserID,
		},
	)
	if err != nil {
		log.Printf("Error posting chirp by '%v' to database: %v", c.UserID, err)
		log.Printf("chirp body: \n\n%v\n", c.Body)
		respondWithError(w, http.StatusInternalServerError, "Error posting chirp: "+err.Error())
		return
	}

	resp := chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}
