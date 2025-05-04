package handler

import (
	"database/sql"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: missing or malformed token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: invalid token")
		return
	}
	chirpIDstr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "404: Chirp not Found")
		return
	}
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "404: Chirp not Found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "500: Error fetching chirp data")
		return
	}
	if chirp.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "403: user is not author of chirp")
		return
	}
	err = cfg.DbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "500: Error deleting chirp")
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, struct{}{})
}
