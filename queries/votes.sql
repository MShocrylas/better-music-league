-- name: CreateVote :one
INSERT INTO votes (
    user_id, submission_id, round_id, song_rank
)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetVote :one
SELECT * FROM votes WHERE vote_id = ?;

-- name: ListVotesBySubmission :many
SELECT * FROM votes WHERE submission_id = ?;
