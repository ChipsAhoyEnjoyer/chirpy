package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/auth"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
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
	pw, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password: "+err.Error())
		return
	}
	u, err := cfg.DbQueries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: pw,
		},
	)
	if err != nil {
		log.Printf("Error posting email '%v' to database: %v\n", req.Email, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error registering user: "+err.Error())
		return
	}
	resp := User{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
	utils.RespondWithJSON(w, http.StatusCreated, resp)
}
