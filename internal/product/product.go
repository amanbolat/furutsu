package product

import "time"

type Product struct {
	ID          string    `sql:"id"`
	Name        string    `sql:"name"`
	Price       int       `sql:"price"`
	Description *string    `sql:"description"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
}
