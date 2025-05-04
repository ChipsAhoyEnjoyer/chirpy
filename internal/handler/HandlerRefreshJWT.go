package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerRefreshJWT(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
		return
	}

	userID, err := cfg.DbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
			return
		}
		log.Printf("Error fetching Refresh Token from DB: %v\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error connecting to DB")
		return
	}

	newJWT, err := auth.MakeJWT(
		userID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Error creating new JWT for user '%v': %v\n", userID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating new Refresh Token")
		return
	}
	utils.RespondWithJSON(
		w,
		http.StatusOK,
		struct {
			Token string `json:"token"`
		}{Token: newJWT},
	)
}
