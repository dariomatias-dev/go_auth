package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type usersController struct{}

func NewUsersController() *usersController {
	return &usersController{}
}

func (uc usersController) Create(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello World!")
}

func (uc usersController) FindOne(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Param("id"))
}

func (uc usersController) FindAll(ctx *gin.Context) {}

func (uc usersController) Update(ctx *gin.Context) {}

func (uc usersController) Delete(ctx *gin.Context) {}
