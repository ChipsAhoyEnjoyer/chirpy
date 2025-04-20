package main

import (
	"sync/atomic"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
}
