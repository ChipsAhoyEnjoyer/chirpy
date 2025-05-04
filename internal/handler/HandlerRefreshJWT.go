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

	user, err := cfg.DbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
			return
		}
		log.Printf("Error fetching Refresh Token from DB: %v\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error connecting to DB")
		return
	}

	// Check if refresh token is expired or revoked
	if time.Now().UTC().After(user.ExpiresAt) || (user.RevokedAt.Valid) {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: Token expired or revoked")
		return
	}
	newJWT, err := auth.MakeJWT(
		user.UserID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Error creating new JWT for user '%v': %v\n", user.UserID, err)
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
