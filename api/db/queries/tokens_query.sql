-- name: CreateTokens :one
INSERT INTO
    "tokens" (
        user_id, access_token, refresh_token
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateTokens :one
UPDATE "tokens"
SET
    access_token = $2,
    refresh_token = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1 RETURNING *;

-- name: DeleteTokens :one
DELETE FROM "tokens" WHERE user_id = $1 RETURNING *;