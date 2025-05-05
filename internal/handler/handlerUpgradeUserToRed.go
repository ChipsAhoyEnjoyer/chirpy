package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerUpgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "500: unable to read request")
		return
	}
	if req.Event != "user.upgraded" {
		utils.RespondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "400: invalid or malformed ID")
		return
	}
	_, err = cfg.DbQueries.UpdateUserToChirpyRed(
		r.Context(),
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "404: user not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "500: database error")
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, struct{}{})

}
