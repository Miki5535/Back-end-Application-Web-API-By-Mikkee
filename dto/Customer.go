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

// Input DTO for login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Output DTO for customer details in the response
type CustomerResponse struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone"`
	FirstName   string `json:"fname"`
	LastName    string `json:"lname"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

// Output DTO for login response
type LoginResponse struct {
	Message  string           `json:"message"`
	Customer CustomerResponse `json:"customer"`
}
