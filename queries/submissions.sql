-- name: CreateSubmission :one
INSERT INTO submissions (
    song_id, user_id, round_id
)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSubmission :one
SELECT * FROM submissions WHERE submission_id = ?;

-- name: ListSubmissionsByRound :many
SELECT * FROM submissions WHERE round_id = ?;

-- Given a round_id, list all submissions with users and their votes for each
-- TODO: This isn't finished
-- name: Foo :many 
SELECT submissions.song_id, submissions.user_id, votes.song_rank
FROM submissions
LEFT JOIN votes ON submissions.submission_id = votes.submission_id
WHERE submissions.round_id = ?;
