-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: UserLogin :one
SELECT * FROM users
WHERE  email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpdateUserCredentials :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING id, created_at, updated_at, email;
