package datastore

import (
	"context"
	"database/sql"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/georgysavva/scany/pgxscan"
	"time"
)

type DbProduct struct {
	ID          string
	Name        string
	Price       int
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductDataStore struct {
	querier pgxscan.Querier
}

func NewProductDataStore(q pgxscan.Querier) *ProductDataStore {
	return &ProductDataStore{querier: q}
}

func (s ProductDataStore) ListProducts(ctx context.Context) ([]product.Product, error) {
	var arr []DbProduct
	err := pgxscan.Select(ctx, s.querier, &arr, `select * from product`)
	if err != nil {
		return nil, err
	}

	products := make([]product.Product, len(arr))
	for i, p := range arr {
		products[i] = product.Product{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description.String,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	return products, nil
}

func (s ProductDataStore) GetProductById(id string, ctx context.Context) (product.Product, error) {
	var p product.Product
	err := pgxscan.Select(ctx, s.querier, &p, `select * from product where id=$1`, id)
	if err != nil {
		return p, err
	}

	return p, nil
}
