// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: email_validations_query.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createEmailValidation = `-- name: CreateEmailValidation :exec
INSERT INTO
    "email_validations" (
        user_id, verification_code, expiration_time
    )
VALUES ($1, $2, $3)
`

type CreateEmailValidationParams struct {
	UserID           uuid.UUID `json:"user_id"`
	VerificationCode string    `json:"verification_code"`
	ExpirationTime   int32     `json:"expiration_time"`
}

func (q *Queries) CreateEmailValidation(ctx context.Context, arg CreateEmailValidationParams) error {
	_, err := q.db.ExecContext(ctx, createEmailValidation, arg.UserID, arg.VerificationCode, arg.ExpirationTime)
	return err
}

const deleteEmailValidation = `-- name: DeleteEmailValidation :exec
DELETE FROM "email_validations" WHERE user_id = $1
`

func (q *Queries) DeleteEmailValidation(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEmailValidation, userID)
	return err
}

const getEmailValidation = `-- name: GetEmailValidation :one
SELECT user_id, verification_code, expiration_time, created_at FROM "email_validations" WHERE verification_code = $1
`

func (q *Queries) GetEmailValidation(ctx context.Context, verificationCode string) (EmailValidations, error) {
	row := q.db.QueryRowContext(ctx, getEmailValidation, verificationCode)
	var i EmailValidations
	err := row.Scan(
		&i.UserID,
		&i.VerificationCode,
		&i.ExpirationTime,
		&i.CreatedAt,
	)
	return i, err
}

const getEmailValidations = `-- name: GetEmailValidations :many
SELECT user_id, verification_code, expiration_time, created_at FROM "email_validations"
`

func (q *Queries) GetEmailValidations(ctx context.Context) ([]EmailValidations, error) {
	rows, err := q.db.QueryContext(ctx, getEmailValidations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EmailValidations
	for rows.Next() {
		var i EmailValidations
		if err := rows.Scan(
			&i.UserID,
			&i.VerificationCode,
			&i.ExpirationTime,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
