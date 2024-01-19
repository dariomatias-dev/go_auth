package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/services"
)

type usersController struct {
	UsersServices services.UsersService
}

func NewUsersController(
	usersServices services.UsersService,
) *usersController {
	return &usersController{
		UsersServices: usersServices,
	}
}

func (uc usersController) Create(ctx *gin.Context) {
	createUser := models.CreateUserModel{}

	if err := ctx.ShouldBindJSON(&createUser); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid body",
				"error":   err.Error(),
			},
		)
		return
	}

	userID := uc.UsersServices.Create(ctx, createUser)

	verificationCode := ""
	for loop := 0; loop < 6; loop++ {
		verificationCode += fmt.Sprint(rand.Intn(10))
	}

	verificationEmailResponse := uc.UsersServices.SendVerificationEmail(
		verificationCode,
		createUser.Name,
		createUser.Email,
	)

	emailValidation := models.EmailValidationModel{
		UserID:           *userID,
		VerificationCode: verificationCode,
		ExpirationTime: time.Now().Add(
			time.Minute * 15,
		).Unix(),
	}

	uc.UsersServices.CreateEmailValidation(
		ctx,
		emailValidation,
	)

	ctx.JSON(
		http.StatusOK,
		verificationEmailResponse,
	)
}

func (uc usersController) FindOne(ctx *gin.Context) {
	userID := ctx.Param("id")

	ID, _ := uuid.Parse(userID)

	user := uc.UsersServices.FindOne(ctx, ID)

	ctx.JSON(
		http.StatusOK,
		user,
	)
}

func (uc usersController) FindAll(ctx *gin.Context) {
	users := uc.UsersServices.FindAll(ctx)

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

	updatedUser := uc.UsersServices.Update(
		ctx,
		ID,
		updateUser,
	)

	ctx.JSON(
		http.StatusOK,
		updatedUser,
	)
}

func (uc usersController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")

	ID, _ := uuid.Parse(userID)

	deletedUser := uc.UsersServices.Delete(ctx, ID)

	ctx.JSON(
		http.StatusOK,
		deletedUser,
	)
}
