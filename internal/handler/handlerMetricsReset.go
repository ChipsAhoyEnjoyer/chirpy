package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/utils"
)

func (cfg *ApiConfig) HandlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("PLATFORM")
	if env != "dev" {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error 403: Forbidden"))
		return
	}
	cfg.FileserverHits.Store(0)
	err := cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error reseting users table: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error reseting users table: "+err.Error())
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
