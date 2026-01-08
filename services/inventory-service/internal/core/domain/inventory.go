package domain

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model

	ProductID uint   `json:"product_id" gorm:"not null;unique"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`

	Qty int `json:"qty" gorm:"not null;default:0"`
}
