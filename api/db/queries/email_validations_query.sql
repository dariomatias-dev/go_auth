-- name: CreateEmailValidation :exec
INSERT INTO
    "email_validations" (
        user_id, verification_code, expiration_time
    )
VALUES ($1, $2, $3);

-- name: GetEmailValidation :one
SELECT * FROM "email_validations" WHERE verification_code = $1;

-- name: GetEmailValidations :many
SELECT * FROM "email_validations";

-- name: DeleteEmailValidation :exec
DELETE FROM "email_validations" WHERE user_id = $1;