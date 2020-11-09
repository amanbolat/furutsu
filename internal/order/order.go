package order

import "time"

type OrderItem struct {
	Id                 string    `json:"id"`
	OrderId            string    `json:"order_id"`
	ProductName        string    `json:"product_name"`
	ProductDescription string    `json:"product_description"`
	Price              int       `json:"price"`
	Amount             int       `json:"amount"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Status string

const (
	StatusPending Status = "pending"
	StatusPaid    Status = "paid"
)

type Order struct {
	Id              string      `json:"id"`
	UserId          string      `json:"user_id"`
	Status          Status      `json:"status"`
	Items           []OrderItem `json:"items"`
	Savings         int         `json:"savings"`
	Total           int         `json:"total"`
	TotalForPayment int         `json:"total_for_payment"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
