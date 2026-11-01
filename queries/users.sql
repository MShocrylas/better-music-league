-- name: CreateUser :one
INSERT INTO users (
    username
)
VALUES (?)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE user_id = ?;

-- name: ListUsers :many
SELECT * FROM users;
