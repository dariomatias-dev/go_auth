package middlewares

import (
	"github.com/gin-gonic/gin"

	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
	"github.com/dariomatias-dev/go_auth/api/services"
)

func VerifyToken(
	ctx *gin.Context,
	authService services.AuthService,
) {
	payload, _, ok := authService.ValidateToken(
		ctx,
		tokentype.AccessToken,
	)
	if !ok {
		return
	}

	ctx.Set("user", *payload)
}
