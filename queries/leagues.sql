-- name: CreateLeague :one
INSERT INTO leagues (
    name
)
VALUES (?)
RETURNING *;

-- name: GetLeague :one
SELECT * FROM leagues WHERE league_id = ?;

-- name: ListLeagues :many
SELECT * FROM leagues;
