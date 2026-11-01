-- name: CreateSong :one
INSERT INTO songs (
    title, artist, url
)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSong :one
SELECT * FROM songs WHERE song_id = ?;

-- name: ListSongs :many
SELECT * FROM songs;
