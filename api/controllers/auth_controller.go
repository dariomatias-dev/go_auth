package controllers

import (
	"net/http"

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

func (ac authController) Login(ctx *gin.Context) {}

func (ac authController) Refresh(ctx *gin.Context) {}

func (ac authController) ValidateEmail(ctx *gin.Context) {
	validateEmail := models.ValidateEmailModel{}

	if err := ctx.ShouldBindJSON(&validateEmail); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid verification code",
				"error":   err,
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		validateEmail.VerificationCode,
	)
}
