package models

type LoginModel struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=10,max=20"`
}
