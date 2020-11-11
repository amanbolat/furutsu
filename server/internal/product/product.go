package product

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
