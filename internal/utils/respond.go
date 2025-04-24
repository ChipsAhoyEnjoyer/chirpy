package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	respBody := struct {
		Err string `json:"error"`
	}{Err: msg}
	dat, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, "Failed to encode error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(payload) // payload should be a Go struct or any JSON-marshalable type
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to encode response: "+err.Error())
		return
	}
	w.Write(dat)
}
