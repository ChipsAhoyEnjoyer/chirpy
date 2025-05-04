package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("Error retreiving token: " + err.Error())
		utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
		return
	}
	if _, err := cfg.DbQueries.UpdateRefreshTokenRevoked(
		r.Context(),
		database.UpdateRefreshTokenRevokedParams{
			UpdatedAt: time.Now().UTC(),
			RevokedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			Token: refreshToken,
		},
	); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusUnauthorized, "401: Unauthorized")
			return
		}
	}
	utils.RespondWithJSON(
		w,
		http.StatusNoContent,
		struct{}{},
	)
}
