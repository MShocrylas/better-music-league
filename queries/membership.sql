-- name: IsRoundMember :one
SELECT EXISTS(
    SELECT 1
    FROM rounds r
    JOIN leagues_users lu ON lu.league_id = r.league_id
    WHERE
        r.round_id = sqlc.arg(round_id)
        AND lu.user_id = sqlc.arg(user_id)
) AS is_member;
