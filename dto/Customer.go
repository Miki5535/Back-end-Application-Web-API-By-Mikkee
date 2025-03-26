package dto

import "time"

type CustomerDTO struct {
	CustomerID  int       `json:"customer_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Password    string    `json:"password"` // Although you might not send password in a DTO
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Carts       []CartDTO `json:"carts"` // List of Carts associated with the customer
}
