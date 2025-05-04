-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshTokenRevoked :one
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $2
WHERE token = $3
RETURNING *;