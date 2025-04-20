package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("PLATFORM")
	if env != "dev" {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Error 403: Forbidden"))
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error reseting users table: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error reseting users table: "+err.Error())
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
