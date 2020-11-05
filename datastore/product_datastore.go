package datastore

import (
	"context"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type ProductDataStore struct {
	tx pgx.Tx
}

func NewProductDataStore(tx pgx.Tx) *ProductDataStore {
	return &ProductDataStore{tx: tx}
}

func (s ProductDataStore) ListProducts(ctx context.Context) ([]product.Product, error) {
	var arr []product.Product
	err := pgxscan.Select(ctx, s.tx, &arr, `select * from product`)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (s ProductDataStore) GetProductById(id string, ctx context.Context) (product.Product, error) {
	var p product.Product
	err := pgxscan.Select(ctx, s.tx, &p, `select * from product where id=$1`, id)
	if err != nil {
		return p, err
	}

	return p, nil
}
