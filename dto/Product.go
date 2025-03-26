package dto

import "time"

type ProductDTO struct {
	ProductID     int       `json:"pid"`
	ProductName   string    `json:"pname"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock"`
	CreatedAt     time.Time `json:"created"`
	UpdatedAt     time.Time `json:"updated"`
}
