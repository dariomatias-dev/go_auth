package controllers

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/services"
)

type authController struct {
	UsersService services.UsersService
	AuthService  services.AuthService
}

func NewAuthController(
	authService services.AuthService,
	usersService services.UsersService,
) *authController {
	return &authController{
		AuthService:  authService,
		UsersService: usersService,
	}
}

func (ac authController) Login(ctx *gin.Context) {
	loginBody := models.LoginModel{}

	if err := ctx.ShouldBindJSON(&loginBody); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"message": "invalid body.",
				"error":   err,
			},
		)
		return
	}

	user := ac.UsersService.FindOneByEmail(
		ctx,
		loginBody.Email,
	)

	if user != nil {
		if !user.ValidEmail {
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"message": "Email not verified.",
				},
			)
			return
		}

		loginAttempts := ac.AuthService.GetLoginAttempts(ctx, user.ID)

		if loginAttempts.Attempts == 10 {
			currentTime := time.Now()
			lastFailedLoginDate := loginAttempts.LastFailedLoginDate.Add(time.Hour * 24)
			hoursLeft := lastFailedLoginDate.Sub(
				currentTime,
			).Hours()

			if hoursLeft > 0 {
				_, value := math.Modf(hoursLeft)
				if value != 0 {
					hoursLeft++
				}

				timeLeft := int(hoursLeft)
				errorMessage := fmt.Sprintf(
					"Your account has been temporarily blocked due to multiple unsuccessful login attempts. Please wait for %d hours before trying again. If issues persist, contact support.",
					timeLeft,
				)

				ctx.JSON(
					http.StatusUnauthorized,
					gin.H{
						"message": errorMessage,
					},
				)
				return
			}

			ac.AuthService.ResetLoginAttempts(ctx, user.ID)
		}

		validPassword := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(loginBody.Password),
		)

		if validPassword == nil {
			tokens := ac.AuthService.GenerateTokens(
				ctx,
				user.ID,
				user.Roles,
			)

			ac.AuthService.UpdateUserTokens(
				ctx,
				user.ID,
				tokens,
			)
			return
		}

		ac.AuthService.IncrementLoginAttemptCounter(
			ctx,
			user.ID,
		)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "invalid email or password.",
		},
	)
}

func (ac authController) Refresh(ctx *gin.Context) {
	payload, tokenString, ok := ac.AuthService.ValidateToken(
		ctx,
		tokentype.RefreshToken,
	)
	if !ok {
		return
	}

	userID, _ := uuid.Parse(payload.ID)

	userTokens := ac.AuthService.GetUserTokens(
		ctx,
		userID,
	)

	if userTokens.RefreshToken != tokenString {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "invalid token",
			},
		)
		return
	}

	tokens := ac.AuthService.GenerateTokens(
		ctx,
		userID,
		payload.Roles,
	)

	ac.AuthService.UpdateUserTokens(
		ctx,
		userID,
		tokens,
	)
}

func (ac authController) ValidateEmail(ctx *gin.Context) {
	validateEmailBody := models.ValidateEmailModel{}

	if err := ctx.ShouldBindJSON(&validateEmailBody); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid verification code",
				"error":   err,
			},
		)
		return
	}

	emailValidation := ac.AuthService.ValidateEmail(
		ctx,
		validateEmailBody.VerificationCode,
	)
	if emailValidation == nil {
		return
	}

	ac.AuthService.DeleteEmailValidation(
		ctx,
		emailValidation.UserID,
	)

	if float64(emailValidation.ExpirationTime) < float64(time.Now().Unix()) {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Expired validation code",
				"error": []float64{
					float64(emailValidation.ExpirationTime),
					float64(time.Now().Unix()),
				},
			},
		)
		return
	}

	ac.AuthService.UpdateUserEmailStatus(
		ctx,
		emailValidation.UserID,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Validated email",
		},
	)
}
