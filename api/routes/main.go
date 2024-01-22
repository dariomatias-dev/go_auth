package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/controllers"
	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	"github.com/dariomatias-dev/go_auth/api/middlewares"
	"github.com/dariomatias-dev/go_auth/api/services"
)

func AppRoutes(
	router *gin.Engine,
	dbQueries *db.Queries,
) *gin.RouterGroup {
	usersService := services.UsersService{
		DbQueries: dbQueries,
	}

	authController := controllers.NewAuthController(
		services.AuthService{
			DbQueries: dbQueries,
		},
		usersService,
	)
	usersController := controllers.NewUsersController(
		usersService,
	)

	router.Use(middlewares.HandleError())

	validUUIDMiddleware := middlewares.ValidUUIDMiddleware

	app := router.Group("")
	{
		auth := app.Group("")
		{
			auth.POST("/login", authController.Login)
			auth.GET("/refresh", authController.Refresh)
			auth.POST("/validate-email", authController.ValidateEmail)
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
