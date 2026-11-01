-- Add a user to a league
-- name: CreateLeagueUser :exec
INSERT INTO leagues_users (
    user_id, league_id, role
)
VALUES (?, ?, ?);

-- name: ListUsersByLeague :many
SELECT users.user_id, users.username, lu.role
FROM leagues_users lu
JOIN users ON lu.user_id = users.user_id
WHERE lu.league_id = ?;

-- name: ListLeaguesByUser :many
SELECT leagues.league_id, leagues.name
FROM leagues_users lu
JOIN leagues ON lu.league_id = leagues.league_id
WHERE lu.user_id = ?;
