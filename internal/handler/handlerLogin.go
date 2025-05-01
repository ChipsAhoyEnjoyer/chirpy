package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	var dur time.Duration
	if 0 < req.ExpiresInSeconds && req.ExpiresInSeconds < 3600 {
		dur = time.Second * time.Duration(req.ExpiresInSeconds)
	} else {
		dur = time.Hour
	}
	tk, err := auth.MakeJWT(
		u.ID,
		cfg.JWTSecret,
		dur,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating session token")
		return
	}
	resp := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
		Token:     tk,
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
