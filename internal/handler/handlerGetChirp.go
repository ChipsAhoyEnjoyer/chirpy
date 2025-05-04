package handler

import (
	"database/sql"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "404: Chirp not Found")
		return
	}
	c, err := cfg.DbQueries.GetChirp(
		r.Context(),
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "404: Chirp not Found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "500: error getting request from database")
		return
	}
	utils.RespondWithJSON(
		w,
		http.StatusOK,
		Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		},
	)
}
