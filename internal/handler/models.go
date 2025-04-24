package handler

import (
	"sync/atomic"
	"time"

	"github.com/ChipsAhoyEnjoyer/chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      database.Queries
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
