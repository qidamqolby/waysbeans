package models

import "time"

type Cart struct {
	ID            int32       `json:"id" gorm:"primary_key: auto_increment"`
	UserID        int32       `json:"user_id"`
	ProductID     int32       `json:"product_id"`
	Product       Product     `json:"product"`
	OrderQty      int32       `json:"orderQuantity"`
	Subtotal      int32       `json:"subtotal"`
	TransactionID int32       `json:"-"`
	Transaction   Transaction `json:"-" gorm:"constraint :OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreateAt      time.Time   `json:"-"`
}
