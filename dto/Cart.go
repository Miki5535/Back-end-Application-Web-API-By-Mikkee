package dto

import "time"

type CartDTO struct {
	CartID     int           `json:"cart_id"`
	CustomerID int           `json:"customer_id"`
	CartName   string        `json:"cart_name"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Items      []CartItemDTO `json:"items"` // List of CartItems
}
