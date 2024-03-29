package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
	"github.com/dariomatias-dev/go_auth/api/models"
)

type AuthService struct {
	DbQueries *db.Queries
}

func (as AuthService) GetToken(
	ctx *gin.Context,
) (string, bool) {
	tokenString, err := GetAuthorizationToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": err.Error(),
				"error":   "token required",
			},
		)
		return "", false
	}

	return tokenString, true
}

func GetAuthorizationToken(
	ctx *gin.Context,
) (string, error) {
	authorization := ctx.GetHeader("Authorization")

	if index := strings.Index(authorization, " "); index == -1 {
		return "", errors.New("invalid token")
	}

	authorizationToken := strings.Split(authorization, " ")
	typeToken := authorizationToken[0]
	tokenString := authorizationToken[1]

	if typeToken != "Bearer" {
		return "", errors.New("invalid token")
	}

	return tokenString, nil
}

func (as AuthService) GenerateTokens(
	ctx *gin.Context,
	userID uuid.UUID,
	userRoles []string,
) *models.Tokens {
	accessToken := generateToken(
		userID,
		userRoles,
		tokentype.AccessToken,
		1,
	)
	refreshToken := generateToken(
		userID,
		userRoles,
		tokentype.RefreshToken,
		7,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
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

func (as AuthService) ValidateToken(
	ctx *gin.Context,
	tokenType string,
) (*models.PayloadModel, string, bool) {
	tokenString, ok := as.GetToken(ctx)
	if !ok {
		return nil, "", false
	}

	payload, ok := as.GetPayload(
		ctx,
		tokenString,
	)
	if !ok {
		return nil, "", false
	}

	userID, _ := uuid.Parse(payload.ID)

	userTokens := as.GetUserTokens(
		ctx,
		userID,
	)
	if userTokens == nil {
		return nil, "", false
	}

	if tokenType != payload.TokenType ||
		(payload.TokenType == tokentype.AccessToken && userTokens.AccessToken != tokenString) ||
		(payload.TokenType == tokentype.RefreshToken && userTokens.RefreshToken != tokenString) {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "invalid token",
			},
		)
		return nil, "", false
	}

	return payload, tokenString, true
}

func (as AuthService) GetUserTokens(
	ctx *gin.Context,
	userID uuid.UUID,
) *db.Tokens {
	userTokens, err := as.DbQueries.GetTokens(ctx, userID)

	if err == sql.ErrNoRows {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"message": "token user not found",
			},
		)
		return nil
	} else if err != nil {
		panic(err)
	}

	return &userTokens
}

func (as AuthService) GetPayload(
	ctx *gin.Context,
	tokenString string,
) (*models.PayloadModel, bool) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return []byte(
				os.Getenv("JWT_SECRET_KEY"),
			), nil
		},
	)

	if err != nil {
		if !token.Valid {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid token",
					"error":   "expired token",
				},
			)
			return nil, false
		}

		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "invalid token",
				"error":   err.Error(),
			},
		)
		return nil, false
	}

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		var roles []string
		userRoles := mapClaims["roles"].([]interface{})

		for _, userRole := range userRoles {
			roles = append(roles, userRole.(string))
		}

		return &models.PayloadModel{
			Exp:       mapClaims["exp"].(float64),
			ID:        mapClaims["id"].(string),
			Roles:     roles,
			TokenType: mapClaims["token_type"].(string),
		}, true
	}

	ctx.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{
			"message": "invalid token",
		},
	)
	return nil, false
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

func (as AuthService) GetLoginAttempts(
	ctx *gin.Context,
	userID uuid.UUID,
) db.GetLoginAttemptRow {
	loginAttempt, err := as.DbQueries.GetLoginAttempt(ctx, userID)
	if err != nil {
		panic(err)
	}

	return loginAttempt
}

func (as AuthService) IncrementLoginAttemptCounter(
	ctx *gin.Context,
	userID uuid.UUID,
) {
	err := as.DbQueries.IncrementLoginAttemptCounter(ctx, userID)
	if err != nil {
		panic(err)
	}
}

func (ac AuthService) ResetLoginAttempts(
	ctx *gin.Context,
	userID uuid.UUID,
) {
	err := ac.DbQueries.ResetLoginAttempts(ctx, userID)
	if err != nil {
		panic(err)
	}
}

func (as AuthService) UpdateUserTokens(
	ctx *gin.Context,
	userID uuid.UUID,
	tokens *models.Tokens,
) {
	UpdateTokensParams := db.UpdateTokensParams{
		UserID:       userID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	err := as.DbQueries.UpdateTokens(
		ctx,
		UpdateTokensParams,
	)
	if err != nil {
		panic(err)
	}
}
