package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/controllers"
	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	"github.com/dariomatias-dev/go_auth/api/middlewares"
)

func AppRoutes(
	router *gin.Engine,
	dbQueries *db.Queries,
) *gin.RouterGroup {
	authController := controllers.NewAuthController()
	usersController := controllers.NewUsersController(
		dbQueries,
	)

	router.Use(middlewares.HandleError())

	validUUIDMiddleware := middlewares.ValidUUIDMiddleware

	app := router.Group("")
	{
		auth := app.Group("")
		{
			auth.POST("/login", authController.Login)
			auth.GET("/refresh", authController.Refresh)
		}

		users := app.Group("")
		{
			users.POST(
				"/user",
				usersController.Create,
			)
			users.GET(
				"/user/:id",
				validUUIDMiddleware,
				usersController.FindOne,
			)
			users.GET(
				"/users",
				usersController.FindAll,
			)
			users.PATCH(
				"/user/:id",
				validUUIDMiddleware,
				usersController.Update,
			)
			users.DELETE(
				"/user/:id",
				validUUIDMiddleware,
				usersController.Delete,
			)
		}
	}

	return app
}
