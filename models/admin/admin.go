package admin

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
