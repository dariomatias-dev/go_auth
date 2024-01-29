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

func (uc usersController) Create(
	ctx *gin.Context,
	userRoles []string,
) {
	createUserBody := models.CreateUserModel{}

	if err := ctx.ShouldBindJSON(&createUserBody); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid body",
				"error":   err.Error(),
			},
		)
		return
	}

	userID := uc.UsersServices.Create(
		ctx,
		createUserBody,
		userRoles,
	)

	verificationCode := ""
	for loop := 0; loop < 6; loop++ {
		verificationCode += fmt.Sprint(rand.Intn(10))
	}

	verificationEmailResponse := uc.UsersServices.SendVerificationEmail(
		verificationCode,
		createUserBody.Name,
		createUserBody.Email,
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
	updateUserBody := models.UpdateModel{}

	ID, _ := uuid.Parse(userID)

	if err := ctx.ShouldBindJSON(&updateUserBody); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid body",
				"error":   err.Error(),
			},
		)
		return
	}

	uc.UsersServices.Update(
		ctx,
		ID,
		updateUserBody,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Updated user",
		},
	)
}

func (uc usersController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")

	ID, _ := uuid.Parse(userID)

	uc.UsersServices.Delete(ctx, ID)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "User deleted",
		},
	)
}
