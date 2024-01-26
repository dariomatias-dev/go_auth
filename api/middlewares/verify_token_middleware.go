package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/services"
)

func VerifyToken(
	ctx *gin.Context,
	authService services.AuthService,
) {

}
