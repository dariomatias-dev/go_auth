package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	tokentype "github.com/dariomatias-dev/go_auth/api/enums/token_type"
	"github.com/dariomatias-dev/go_auth/api/services"
)

func VerifyToken(
	ctx *gin.Context,
	authService services.AuthService,
) {
	token, ok := authService.GetToken(ctx)
	if !ok {
		return
	}

	payload, ok := authService.GetPayload(
		ctx,
		*token,
	)
	if !ok {
		return
	}

	if mapClaims, ok := payload.Claims.(jwt.MapClaims); ok || payload.Valid {
		if mapClaims["token_type"] != tokentype.AccessToken {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "invalid token",
				},
			)
			return
		}

		ctx.Set("user", payload)
	}
}
