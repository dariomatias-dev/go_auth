-- name: CreateTokens :exec
INSERT INTO "tokens" (user_id) VALUES ($1);

-- name: GetTokens :one
SELECT * FROM "tokens" WHERE user_id = $1;

-- name: GetAllTokens :many
SELECT * FROM "tokens";

-- name: UpdateTokens :exec
UPDATE "tokens"
SET
    access_token = $2,
    refresh_token = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1;

-- name: DeleteTokens :exec
DELETE FROM "tokens" WHERE user_id = $1;