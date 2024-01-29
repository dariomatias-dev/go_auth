package routes

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

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
				"/user-admin",
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

	ctx := context.Background()

	adminEmail := os.Getenv("ADMIN_EMAIL")

	_, err := dbQueries.GetUserByEmail(
		ctx,
		adminEmail,
	)

	if err == sql.ErrNoRows {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(os.Getenv("ADMIN_PASSWORD")),
			10,
		)
		if err != nil {
			log.Fatal(err)
		}

		createUserParams := db.CreateUserParams{
			Name:     "Administrator",
			Age:      18,
			Email:    adminEmail,
			Password: string(encryptedPassword),
			Roles: []string{
				usertype.Admin,
			},
		}

		dbQueries.CreateUser(
			ctx,
			createUserParams,
		)
	} else if err != nil {
		panic(err)
	}

	return app
}
