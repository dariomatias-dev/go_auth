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
	authService := services.AuthService{
		DbQueries: dbQueries,
	}
	usersService := services.UsersService{
		DbQueries: dbQueries,
	}

	authController := controllers.NewAuthController(
		authService,
		usersService,
	)
	usersController := controllers.NewUsersController(
		usersService,
	)

	router.Use(middlewares.HandleError())

	validUUIDMiddleware := middlewares.ValidUUIDMiddleware
	verifyTokenMiddleware := func(ctx *gin.Context) {
		middlewares.VerifyToken(
			ctx,
			authService,
		)
	}

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
				verifyTokenMiddleware,
				validUUIDMiddleware,
				usersController.FindOne,
			)
			users.GET(
				"/users",
				verifyTokenMiddleware,
				usersController.FindAll,
			)
			users.PATCH(
				"/user/:id",
				verifyTokenMiddleware,
				validUUIDMiddleware,
				usersController.Update,
			)
			users.DELETE(
				"/user/:id",
				verifyTokenMiddleware,
				validUUIDMiddleware,
				usersController.Delete,
			)
		}
	}

	return app
}
