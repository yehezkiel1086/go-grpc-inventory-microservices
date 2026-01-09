package domain

import "gorm.io/gorm"

type OrderStatus uint8

const (
	Unpaid OrderStatus = 0
	Paid OrderStatus = 1
	Shipped OrderStatus = 2
	Delivered OrderStatus = 3
)

type Order struct {
	gorm.Model

	UserID uint `json:"user_id" gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`

	ProductID uint `json:"product_id" gorm:"not null"`
	Product Product `gorm:"foreignKey:ProductID"`

	Qty uint `json:"qty" gorm:"not null"`
	TotalPrice float64 `json:"total_price" gorm:"not null"`
	Status OrderStatus `json:"status" gorm:"not null;default:0"`
}
