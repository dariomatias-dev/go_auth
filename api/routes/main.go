package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/controllers"
	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
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
	identityVerifierMiddleware := middlewares.IdentityVerifierMiddleware
	roleCheckMiddleware := middlewares.RoleCheckMiddleware
	adminCheckMiddleware := func(ctx *gin.Context) {
		roles := []string{
			usertype.Admin,
		}

		roleCheckMiddleware(
			ctx,
			roles,
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
				adminCheckMiddleware,
				func(ctx *gin.Context) {
					usersController.Create(
						ctx,
						[]string{
							usertype.Admin,
						},
					)
				},
			)
			users.POST(
				"/user",
				func(ctx *gin.Context) {
					usersController.Create(
						ctx,
						[]string{
							usertype.User,
						},
					)
				},
			)
			users.GET(
				"/user/:id",
				verifyTokenMiddleware,
				validUUIDMiddleware,
				identityVerifierMiddleware,
				usersController.FindOne,
			)
			users.GET(
				"/users",
				verifyTokenMiddleware,
				adminCheckMiddleware,
				usersController.FindAll,
			)
			users.PATCH(
				"/user/:id",
				verifyTokenMiddleware,
				validUUIDMiddleware,
				identityVerifierMiddleware,
				usersController.Update,
			)
			users.DELETE(
				"/user/:id",
				verifyTokenMiddleware,
				validUUIDMiddleware,
				identityVerifierMiddleware,
				usersController.Delete,
			)
		}
	}

	return app
}
