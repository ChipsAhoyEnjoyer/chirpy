package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
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
	resp := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
