package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding response")
		return
	}

	if len(req.Body) > 140 {
		log.Println("Error with user request: Too many characters")
		utils.RespondWithError(w, http.StatusBadRequest, "Request too long; Max 140 char request.")
		return
	}

	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(tk, cfg.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
		return
	}

	cleanedChirp := utils.ProfanityFilter(req.Body)
	c, err := cfg.DbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedChirp,
			UserID: userID,
		},
	)
	if err != nil {
		log.Printf("Error posting chirp by '%v' to database: %v", c.UserID, err)
		log.Printf("chirp body: \n\n%v\n", c.Body)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error posting chirp")
		return
	}

	resp := Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
	utils.RespondWithJSON(w, http.StatusCreated, resp)
}
