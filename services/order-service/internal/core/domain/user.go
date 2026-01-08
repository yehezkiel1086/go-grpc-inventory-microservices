package domain

import "gorm.io/gorm"

type Role uint16

const (
	AdminRole Role = 5150
	UserRole  Role = 2001
)

type User struct {
	gorm.Model

	Email    string `json:"email" gorm:"size:255;not null;unique"`
	Password string `json:"password" gorm:"size:255;not null"`
	Name     string `json:"name" gorm:"size:255;not null"`
	Role     Role   `json:"role" gorm:"default:2001;not null"`
}

type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type UserResponse struct {
	gorm.Model

	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
}
