package models

type ValidateEmailModel struct {
	VerificationCode string `json:"verification_code" binding:"required,min=6,max=6"`
}
