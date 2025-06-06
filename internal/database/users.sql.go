// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    false
)
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

type CreateUserRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const updateUserCredentials = `-- name: UpdateUserCredentials :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type UpdateUserCredentialsParams struct {
	Email          string
	HashedPassword string
	UpdatedAt      time.Time
	ID             uuid.UUID
}

type UpdateUserCredentialsRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateUserCredentials(ctx context.Context, arg UpdateUserCredentialsParams) (UpdateUserCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserCredentials,
		arg.Email,
		arg.HashedPassword,
		arg.UpdatedAt,
		arg.ID,
	)
	var i UpdateUserCredentialsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserToChirpyRed = `-- name: UpdateUserToChirpyRed :one
UPDATE users
set is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type UpdateUserToChirpyRedRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateUserToChirpyRed(ctx context.Context, id uuid.UUID) (UpdateUserToChirpyRedRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserToChirpyRed, id)
	var i UpdateUserToChirpyRedRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const userLogin = `-- name: UserLogin :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE  email = $1
`

func (q *Queries) UserLogin(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, userLogin, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
