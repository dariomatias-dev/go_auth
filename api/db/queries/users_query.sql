-- name: CreateUser :one
INSERT INTO
    "users" (
        name, age, email, password, roles
    )
VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: GetUser :one
SELECT * FROM "users" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE email = $1;

-- name: GetUsers :many
SELECT
    id,
    name,
    age,
    email,
    roles,
    created_at,
    updated_at
FROM "users";

-- name: UpdateUser :one
UPDATE "users"
SET
    name = COALESCE(sqlc.narg ('name'), name),
    age = COALESCE(sqlc.narg ('age'), age),
    email = COALESCE(sqlc.narg ('email'), email),
    valid_email = COALESCE(
        sqlc.narg ('valid_email'), valid_email
    ),
    password = COALESCE(
        sqlc.narg ('password'), password
    ),
    roles = COALESCE(sqlc.narg ('roles'), roles),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING id,
    name,
    age,
    email,
    roles,
    created_at,
    updated_at;

-- name: DeleteUser :one
DELETE FROM "users"
WHERE
    id = $1 RETURNING id,
    name,
    age,
    email,
    roles,
    created_at,
    updated_at;