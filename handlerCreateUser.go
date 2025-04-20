package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := struct {
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding response")
	}
	log.Println(req.Email)
	resp := struct {
		OK string `json:"ok"`
	}{OK: "ok"}
	respondWithJSON(w, http.StatusOK, resp)
}
