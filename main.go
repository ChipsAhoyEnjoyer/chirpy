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
		w.WriteHeader(http.StatusInternalServerError)
		dat, err := json.Marshal(&struct {
			Error string `json:"error"`
		}{Error: "Something went wrong"})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		log.Println("Error with user request: Too many characters")
		w.WriteHeader(http.StatusBadRequest)
		dat, err := json.Marshal(&struct {
			Error string `json:"error"`
		}{Error: "Chirp is too long"})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(dat)
		return
	}
	fmt.Println(params.Body)
	type respBody struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}
	dat, err := json.Marshal(respBody{Valid: true})
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		dat, err = json.Marshal(&respBody{Error: "Something went wrong"})
		if err != nil {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(dat)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
