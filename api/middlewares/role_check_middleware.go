package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/models"
)

func RoleCheckMiddleware(
	ctx *gin.Context,
	roles []string,
) {
	payload, _ := ctx.Get("user")
	user := payload.(models.PayloadModel)

	canAccessRoute := true

	for _, userRole := range user.Roles {
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
