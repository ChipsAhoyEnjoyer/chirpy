package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	port    = "8080"
	rootDir = "."
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	},
	)
}

func (cfg *apiConfig) handlerMetricsCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(
		"<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>",
		cfg.fileserverHits.Load(),
	)
	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handlerServerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Println("Error with user request: Too many characters")
		respondWithError(w, http.StatusBadRequest, "Request too long; Max 140 char request.")
		return
	}
	payload := struct {
		CleanedBody string `json:"cleaned_body,omitempty"`
	}{CleanedBody: profanityFilter(params.Body)}
	respondWithJSON(w, http.StatusOK, payload)
}

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

func main() {
	cfg := apiConfig{fileserverHits: atomic.Int32{}}
	// New router
	mux := http.NewServeMux()

	// Endpoint handlers
	fileServer := http.FileServer(http.Dir(rootDir))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServer))) // Requests should start with /app/ to avoid
	mux.HandleFunc("GET /api/healthz", handlerServerReady)                              // conflicts with other endpoints
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetricsCount)
	mux.HandleFunc("POST /admin/reset", cfg.handlerMetricsReset)

	server := &http.Server{ // Server configurations
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving '%v' on port: %v\n", rootDir, port)
	log.Fatal(server.ListenAndServe())
}
