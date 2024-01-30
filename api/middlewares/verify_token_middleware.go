package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
	"github.com/dariomatias-dev/go_auth/api/services"
)

func VerifyToken(
	ctx *gin.Context,
	authService services.AuthService,
) {
	tokenString, ok := authService.GetToken(ctx)
	if !ok {
		return
	}

	payload, ok := authService.GetPayload(
		ctx,
		tokenString,
	)
	if !ok {
		return
	}

	if payload.TokenType != tokentype.AccessToken {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "invalid token",
			},
		)
		return
	}

	ctx.Set("user", *payload)
}
