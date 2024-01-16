package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/controllers"
)

func AppRoutes(router *gin.Engine) *gin.RouterGroup {
	authController := controllers.NewAuthController()
	usersController := controllers.NewUsersController()

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
				usersController.Create,
			)
			users.GET(
				"/user",
				usersController.Create,
			)
			users.PATCH(
				"/user/:id",
				usersController.Create,
			)
			users.DELETE(
				"/user/:id",
				usersController.Create,
			)
		}
	}

	return app
}
