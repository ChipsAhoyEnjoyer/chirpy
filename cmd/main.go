package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/ChipsAhoyEnjoyer/chirpy/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	port    = "8080"
	rootDir = "."
)

func main() {
	err := godotenv.Load() // load up .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL) // load up postgres connection
	if err != nil {
		log.Fatal("Error opening database connection")
	}
	cfg := handler.ApiConfig{
		FileserverHits: atomic.Int32{},
		DbQueries:      *database.New(db),
	}

	// New router
	mux := http.NewServeMux()

	// Endpoint handlers
	fileServer := http.FileServer(http.Dir(rootDir))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(fileServer))) // Requests should start with /app/ to avoid
	mux.HandleFunc("GET /api/healthz", handler.HandlerServerReady)                      // conflicts with other endpoints
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.HandlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetricsCount)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerMetricsReset)

	server := &http.Server{ // Server configurations
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving '%v' on port: %v\n", rootDir, port)
	log.Fatal(server.ListenAndServe())
}
