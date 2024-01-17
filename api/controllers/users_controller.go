package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/utils"
)

type usersController struct {
	DbQueries *db.Queries
}

func NewUsersController(dbQueries *db.Queries) *usersController {
	return &usersController{
		DbQueries: dbQueries,
	}
}

func (uc usersController) Create(ctx *gin.Context) {
	createUser := models.CreateUserModel{}

	if err := ctx.ShouldBindJSON(&createUser); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(createUser.Password),
		10,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	// Create User
	userArg := db.CreateUserParams{
		Name:     createUser.Name,
		Age:      createUser.Age,
		Email:    createUser.Email,
		Password: string(encryptedPassword),
		Roles: []string{
			usertype.User,
		},
	}

	createdUser, err := uc.DbQueries.CreateUser(ctx, userArg)

	if err != nil {
		panic(err)
	}

	// Create Tokens
	tokensArg := db.CreateTokensParams{
		UserID:       createdUser.ID,
		AccessToken:  "",
		RefreshToken: "",
	}

	_, err = uc.DbQueries.CreateTokens(ctx, tokensArg)

	if err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusOK,
		createdUser,
	)
}

func (uc usersController) FindOne(ctx *gin.Context) {
	userID := ctx.Param("id")

	ID, _ := uuid.Parse(userID)

	user, err := uc.DbQueries.GetUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusOK,
		user,
	)
}

func (uc usersController) FindAll(ctx *gin.Context) {
	users, err := uc.DbQueries.GetUsers(ctx)
	if err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusOK,
		users,
	)
}

func (uc usersController) Update(ctx *gin.Context) {
	userID := ctx.Param("id")
	updateUser := models.UpdateModel{}

	ID, _ := uuid.Parse(userID)

	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid body",
				"error":   err.Error(),
			},
		)
		return
	}

	getValue := utils.GetValue{}

	updateUserParams := db.UpdateUserParams{
		ID: ID,
		Name: getValue.String(updateUser.Name),
		Age: getValue.Int32(updateUser.Age),
		Email: getValue.String(updateUser.Email),
	}

	updatedUser, err := uc.DbQueries.UpdateUser(ctx, updateUserParams)
	if err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusOK,
		updatedUser,
	)
}

func (uc usersController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")

	ID, _ := uuid.Parse(userID)

	_, err := uc.DbQueries.DeleteTokens(ctx, ID)
	if err != nil {
		panic(err)
	}

	deletedUser, err := uc.DbQueries.DeleteUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusOK,
		deletedUser,
	)
}
