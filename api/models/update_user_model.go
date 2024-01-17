package models

type UpdateModel struct {
	Name     *string `json:"name" binding:"omitempty,min=3,max=128"`
	Age      *int32  `json:"age" binding:"omitempty,min=18,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
}
