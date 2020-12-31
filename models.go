package platform_exercise

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"-"`
	ID        string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=0"`
}

type UpdateUserRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gt=0"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}
