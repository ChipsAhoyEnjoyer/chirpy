package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerUpdateUserCred(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: missing or incorrect token")
		return
	}
	userID, err := auth.ValidateJWT(tk, cfg.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401: invalid token")
		return
	}
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "500: malformed email/password")
		return
	}
	hashedPW, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "400: invalid password")
		return
	}
	row, err := cfg.DbQueries.UpdateUserCredentials(
		r.Context(),
		database.UpdateUserCredentialsParams{
			Email:          req.Email,
			HashedPassword: hashedPW,
			UpdatedAt:      time.Now().UTC(),
			ID:             userID,
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "500: error updating user credentials")
		return
	}
	utils.RespondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:        row.ID,
			UpdatedAt: row.UpdatedAt,
			CreatedAt: row.CreatedAt,
			Email:     row.Email,
		},
	)
}
