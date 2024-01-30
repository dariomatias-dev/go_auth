package models

import "github.com/google/uuid"

type EmailValidationModel struct {
	UserID           uuid.UUID
	VerificationCode string
	ExpirationTime   int64
}
