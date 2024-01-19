package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/services"
)

type authController struct {
	AuthService services.AuthService
}

func NewAuthController(
	authService services.AuthService,
) *authController {
	return &authController{
		AuthService: authService,
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

	ctx.JSON(
		http.StatusOK,
		loginBody,
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
