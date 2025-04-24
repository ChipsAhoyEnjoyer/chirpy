package handler

import (
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	c, err := cfg.DbQueries.GetChirp(
		r.Context(),
		id,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting request from database: "+err.Error())
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
