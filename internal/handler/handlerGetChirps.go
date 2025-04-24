package handler

import (
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cc, err := cfg.DbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching chirps from database: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error posting chirp: "+err.Error())
		return
	}

	resp := []Chirp{}
	for _, c := range cc {
		resp = append(
			resp,
			Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			},
		)
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
