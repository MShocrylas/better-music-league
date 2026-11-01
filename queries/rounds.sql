-- name: CreateRound :one
INSERT INTO rounds (
    league_id, theme, position, submission_deadline, voting_deadline
)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRound :one
SELECT * FROM rounds WHERE round_id = ?;

-- name: ListRoundsByLeague :many
SELECT * FROM rounds WHERE league_id = ?;
