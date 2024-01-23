-- name: CreateLoginAttempts :exec
INSERT INTO "login_attempts" (user_id) VALUES ($1);

-- name: GetLoginAttempt :one
SELECT * FROM "login_attempts" WHERE user_id = $1;

-- name: GetLoginAttempts :many
SELECT * FROM "login_attempts";

-- name: UpdateLoginAttempts :exec
UPDATE "login_attempts"
SET
    attempts = COALESCE(
        sqlc.narg ('attempts'), attempts
    ),
    last_failed_login_date = CURRENT_TIMESTAMP
WHERE
    user_id = $1;

-- name: DeleteLoginAttempt :exec
DELETE FROM "login_attempts" WHERE user_id = $1;