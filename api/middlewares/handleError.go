package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func HandleError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if pqErr, ok := err.(*pq.Error); ok {
					handleDatabaseError(ctx, pqErr)
				} else if err == sql.ErrNoRows {
					ctx.AbortWithStatusJSON(
						http.StatusNotFound,
						nil,
					)
					return
				}

				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"message": "Internal server error",
						"error": err,
					},
				)
				return
			}
		}()

		ctx.Next()
	}
}

func handleDatabaseError(ctx *gin.Context, pqErr *pq.Error) {
	switch pqErr.Code {
	case "23505":
		fieldName := extractFieldName(pqErr.Message)
		errorMessage := fmt.Sprintf(
			"Vallue of the %s already exists",
			fieldName,
		)

		ctx.AbortWithStatusJSON(
			http.StatusConflict,
			gin.H{
				"message":   errorMessage,
				"fieldName": fieldName,
			},
		)
		return
	}
}

func extractFieldName(message string) string {
	underlineChar := "_"
	fieldName := ""

	startIndex := strings.Index(message, underlineChar)
	if startIndex == -1 {
		return fieldName
	}

	endIndex := strings.LastIndex(message, underlineChar)
	if endIndex == -1 {
		return fieldName
	}

	fieldName = message[startIndex+1 : endIndex]

	return fieldName
}
