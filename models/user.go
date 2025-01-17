package models

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`

	UserName   string `json:"user_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Status     string `json:"status" validate:"required,oneof=A I T"`
	Department string `json:"department" validate:"required"`
}
