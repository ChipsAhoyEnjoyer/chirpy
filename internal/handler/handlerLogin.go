package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

const (
	oneHour   = time.Hour
	sixtyDays = time.Hour * 24 * 60
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding response: "+err.Error())
		return
	}
	u, err := cfg.DbQueries.UserLogin(
		r.Context(),
		req.Email,
	)
	if err != nil {
		log.Printf("Error logging in user '%v': %v\n", req.Email, err)
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if err = auth.CheckPasswordHash(u.HashedPassword, req.Password); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	dur := time.Hour
	tk, err := auth.MakeJWT(
		u.ID,
		cfg.JWTSecret,
		dur,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating session token")
		return
	}
	refreshTk, _ := auth.MakeRefreshToken()
	utcNow := time.Now().UTC()
	_, err = cfg.DbQueries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshTk,
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			ExpiresAt: utcNow.Add(sixtyDays),
			UserID:    u.ID,
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating session token")
		return
	}
	resp := User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Token:        tk,
		RefreshToken: refreshTk,
		IsChirpyRed:  u.IsChirpyRed,
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
