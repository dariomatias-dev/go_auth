package controllers

import "github.com/gin-gonic/gin"

type authController struct{}

func NewAuthController() *authController {
	return &authController{}
}

func (ac authController) Login(ctx *gin.Context) {}

func (ac authController) Refresh(ctx *gin.Context) {}
