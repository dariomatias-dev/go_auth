package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
	"github.com/dariomatias-dev/go_auth/api/models"
)

func IdentityVerifierMiddleware(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, _ := ctx.Get("user")
	payload, _ := user.(models.PayloadModel)

	isAdmin := false

	for _, userRole := range payload.Roles {
		if userRole == usertype.Admin {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		if userID != payload.ID {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "You do not have permission to access another user's data",
				},
			)
		}
	}
}
