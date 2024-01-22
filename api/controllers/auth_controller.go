package controllers

import (
	"net/http"
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

		validPassword := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(loginBody.Password),
		)

		if validPassword == nil {
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"message": "Authenticated",
				},
			)
			return
		}
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Invalid email or password",
		},
	)
}

func (ac authController) Refresh(ctx *gin.Context) {}

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
