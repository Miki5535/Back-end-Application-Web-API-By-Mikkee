package dto

import "time"

type CartItemDTO struct {
	CartItemID int        `json:"cart_item_id"`
	CartID     int        `json:"cart_id"`
	ProductID  int        `json:"product_id"`
	Quantity   int        `json:"quantity"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Product    ProductDTO `json:"product"` // Include Product DTO if needed
}

type CartItemResponseDTO struct {
	ProductID    int     `json:"pid"`
	ProductName  string  `json:"pname"`
	Description  string  `json:"description,omitempty"`
	Quantity     int     `json:"quantity"`
	PricePerUnit float64 `json:"price"`
	TotalPrice   float64 `json:"total"`
}

type CartResponseDTO struct {
	CartID int                   `json:"cid"`
	Name   string                `json:"cname"`
	Items  []CartItemResponseDTO `json:"items"`
}
