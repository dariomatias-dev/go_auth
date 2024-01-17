package controllers

import (
	"github.com/gin-gonic/gin"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
)

type usersController struct {
	DbQueries *db.Queries
}

func NewUsersController(dbQueries *db.Queries) *usersController {
	return &usersController{
		DbQueries: dbQueries,
	}
}

func (uc usersController) Create(ctx *gin.Context) {}

func (uc usersController) FindOne(ctx *gin.Context) {}

func (uc usersController) FindAll(ctx *gin.Context) {}

func (uc usersController) Update(ctx *gin.Context) {}

func (uc usersController) Delete(ctx *gin.Context) {}
