package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	respBody := struct {
		Err string `json:"error"`
	}{Err: msg}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	dat, _ := json.Marshal(payload) // payload should be an encoded struct
	w.Write(dat)
}
