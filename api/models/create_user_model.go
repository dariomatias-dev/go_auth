package models

type CreateUserModel struct {
	Name     string `json:"name" binding:"required,min=3,max=128"`
	Age      int32    `json:"age" binding:"required,min=18,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=10,max=20,alphanum"`
}
