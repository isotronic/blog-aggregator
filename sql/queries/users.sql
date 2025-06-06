-- name: CreateUser :one
INSERT INTO users (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: ResetUsers :exec
DELETE FROM users;