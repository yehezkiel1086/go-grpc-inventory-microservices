package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string  `json:"name" gorm:"size:50;not null;unique"`
	Price       float64 `json:"price" gorm:"not null"`
}

type CreateProductReq struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Qty int `json:"qty" binding:"required"`
}

type ProductRes struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty int `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
