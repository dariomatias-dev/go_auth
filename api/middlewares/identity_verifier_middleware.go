package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
)

func IdentityVerifierMiddleware(ctx *gin.Context) {
	userID := ctx.Param("id")

	payload, _ := ctx.Get("user")

	mapClaims, _ := payload.(*jwt.Token).Claims.(jwt.MapClaims)

	userRoles := mapClaims["roles"].([]interface{})

	isAdmin := true

	for _, userRole := range userRoles {
		if userRole == usertype.Admin {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		if userID != mapClaims["id"] {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "You do not have permission to access another user's data",
				},
			)
		}
	}
}
