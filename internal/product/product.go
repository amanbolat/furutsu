package product

import "time"

type Product struct {
	ID          string
	Name        string
	Price       int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
