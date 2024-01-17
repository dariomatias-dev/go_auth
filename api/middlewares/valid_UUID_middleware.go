package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidUUIDMiddleware(ctx *gin.Context) {
	value := ctx.Param("id")

	uuid, err := uuid.Parse(value)

	if err != nil || uuid.Version() != 4 {
		errorMessage := fmt.Sprintf("Invalid UUID: %s", value)

		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": errorMessage,
			},
		)
		return
	}
}
