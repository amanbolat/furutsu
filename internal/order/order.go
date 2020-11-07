package order

import "time"

type OrderItem struct {
	Id                 string
	OrderId            string
	ProductName        string
	ProductDescription string
	Price              int
	Amount             int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Status string

const (
	StatusPending Status = "pending"
	StatusPaid    Status = "paid"
)

type Order struct {
	Id              string
	UserId          string
	Status          Status
	Items           []OrderItem
	Savings         int
	Total           int
	TotalForPayment int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
