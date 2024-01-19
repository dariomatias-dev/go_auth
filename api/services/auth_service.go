package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
)

type AuthService struct {
	DbQueries *db.Queries
}

func (as AuthService) Login() {}

func (as AuthService) Refresh() {}

func (as AuthService) ValidateEmail(
	ctx *gin.Context,
	validationCode string,
) *db.EmailValidations {
	emailValidation, err := as.DbQueries.GetEmailValidation(
		ctx,
		validationCode,
	)
	if err == sql.ErrNoRows {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"message": "Verification code not found",
			},
		)
		return nil
	} else if err != nil {
		panic(err)
	}

	return &emailValidation
}

func (as AuthService) DeleteEmailValidation(
	ctx *gin.Context,
	emailValidationID uuid.UUID,
) {
	err := as.DbQueries.DeleteEmailValidation(ctx, emailValidationID)
	if err != nil {
		panic(err)
	}
}

func (as AuthService) UpdateUserEmailStatus(
	ctx *gin.Context,
	userID uuid.UUID,
) {
	updateUserParams := db.UpdateUserParams{
		ID: userID,
		ValidEmail: sql.NullBool{
			Bool: true,
			Valid: true,
		},
	}

	as.DbQueries.UpdateUser(ctx, updateUserParams)
}
