package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name        string  `json:"name" gorm:"size:50;not null;unique"`
	Price       float64 `json:"price" gorm:"not null"`
}
