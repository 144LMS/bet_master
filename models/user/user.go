package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteUserRequest struct {
	ID int `json:"id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
