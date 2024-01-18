package services

import db "github.com/dariomatias-dev/go_auth/api/db/sqlc"

type AuthService struct {
	DbQueries *db.Queries
}

func (as AuthService) Login() {}

func (as AuthService) Refresh() {}

func (as AuthService) ValidateEmail(
	validationCode string,
) {}
