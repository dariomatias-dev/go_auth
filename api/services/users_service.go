package services

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

type UsersService struct {
	DbQueries *db.Queries
}

func (us UsersService) Create(
	ctx *gin.Context,
	createUser models.CreateUserModel,
) *db.CreateUserRow {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(createUser.Password),
		10,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return nil
	}

	// Create User
	createUserParams := db.CreateUserParams{
		Name:     createUser.Name,
		Age:      createUser.Age,
		Email:    createUser.Email,
		Password: string(encryptedPassword),
		Roles: []string{
			usertype.User,
		},
	}

	createdUser, err := us.DbQueries.CreateUser(ctx, createUserParams)

	if err != nil {
		panic(err)
	}

	// Create Tokens
	createTokensParams := db.CreateTokensParams{
		UserID:       createdUser.ID,
		AccessToken:  "",
		RefreshToken: "",
	}

	_, err = us.DbQueries.CreateTokens(ctx, createTokensParams)

	if err != nil {
		panic(err)
	}

	return &createdUser
}

func (us UsersService) FindOne(
	ctx *gin.Context,
	ID uuid.UUID,
) *db.Users {
	user, err := us.DbQueries.GetUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	return &user
}

func (us UsersService) FindAll(
	ctx *gin.Context,
) *[]db.GetUsersRow {
	users, err := us.DbQueries.GetUsers(ctx)
	if err != nil {
		panic(err)
	}

	return &users
}

func (us UsersService) Update(
	ctx *gin.Context,
	ID uuid.UUID,
	updateUser models.UpdateModel,
) *db.UpdateUserRow {
	getValue := utils.GetValue{}

	password := updateUser.Password

	if password != nil {
		value, err := bcrypt.GenerateFromPassword(
			[]byte(*password),
			10,
		)
		if err != nil {
			panic(err)
		}

		encryptedPassword := string(value)
		password = &encryptedPassword
	}

	updateUserParams := db.UpdateUserParams{
		ID:       ID,
		Name:     getValue.String(updateUser.Name),
		Age:      getValue.Int32(updateUser.Age),
		Email:    getValue.String(updateUser.Email),
		Password: getValue.String(password),
	}

	updatedUser, err := us.DbQueries.UpdateUser(ctx, updateUserParams)
	if err != nil {
		panic(err)
	}

	return &updatedUser
}

func (us UsersService) Delete(
	ctx *gin.Context,
	ID uuid.UUID,
) *db.DeleteUserRow {
	_, err := us.DbQueries.DeleteTokens(ctx, ID)
	if err != nil {
		panic(err)
	}

	deletedUser, err := us.DbQueries.DeleteUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	return &deletedUser
}
