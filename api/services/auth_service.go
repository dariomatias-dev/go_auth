package services

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
)

type AuthService struct {
	DbQueries *db.Queries
}

func (as AuthService) Login() {}

func (as AuthService) Refresh() {}

func (as AuthService) GenerateTokens(
	ctx *gin.Context,
	userID uuid.UUID,
	userRoles []string,
) *[]string {
	access_token := generateToken(
		userID,
		userRoles,
		tokentype.AccessToken,
		1,
	)
	refresh_token := generateToken(
		userID,
		userRoles,
		tokentype.RefreshToken,
		7,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"access_token":  access_token,
			"refresh_token": refresh_token,
		},
	)

	return &[]string{
		access_token,
		refresh_token,
	}
}

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
			Bool:  true,
			Valid: true,
		},
	}

	as.DbQueries.UpdateUser(ctx, updateUserParams)
}

func generateToken(
	userID uuid.UUID,
	userRoles []string,
	tokenType string,
	daysToExpire int,
) string {
	payload := jwt.MapClaims{
		"id":         userID,
		"roles":      userRoles,
		"token_type": tokenType,
		"exp": time.Now().Add(
			time.Hour * 24 * time.Duration(daysToExpire),
		).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		payload,
	)

	tokenString, err := token.SignedString(
		[]byte(
			os.Getenv("JWT_SECRET_KEY"),
		),
	)
	if err != nil {
		return ""
	}

	return tokenString
}
