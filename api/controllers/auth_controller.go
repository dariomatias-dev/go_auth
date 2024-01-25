package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

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
				"message": "Invalid body",
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
		if !user.ValidEmail.Bool {
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"message": "Email not verified",
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
				http.StatusOK,
				gin.H{
					"message": errorMessage,
				},
			)
			return
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
			"message": "Invalid email or password",
		},
	)
}

func (ac authController) Refresh(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")

	if index := strings.Index(authorization, " "); index == -1 {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid token",
			},
		)
		return
	}

	authorizationToken := strings.Split(authorization, " ")
	typeToken := authorizationToken[0]
	token := authorizationToken[1]

	if typeToken != "Bearer" {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid token",
			},
		)
		return
	}

	ctx.AbortWithStatusJSON(
		http.StatusOK,
		token,
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
