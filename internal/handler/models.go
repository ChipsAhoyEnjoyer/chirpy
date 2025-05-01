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
	JWTSecret      string
}
type Chirp struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	Body      string    `json:"body,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
}

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	Email     string    `json:"email,omitempty"`
	Token     string    `json:"token,omitempty"`
}
