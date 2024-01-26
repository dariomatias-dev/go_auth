package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RoleCheckMiddleware(
	ctx *gin.Context,
	roles []string,
) {
	payload, _ := ctx.Get("user")
	mapClaims, _ := payload.(*jwt.Token).Claims.(jwt.MapClaims)
	userRoles := mapClaims["roles"].([]interface{})

	canAccessRoute := true

	for _, userRole := range userRoles {
		for _, role := range roles {
			if userRole == role {
				canAccessRoute = true
			}
		}
	}

	if !canAccessRoute {
		ctx.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"message": "you do not have the necessary permissions to access this route",
			},
		)
	}
}
