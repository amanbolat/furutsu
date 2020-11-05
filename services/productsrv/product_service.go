package productsrv

import (
	"context"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/jackc/pgx/v4"
)

type Service struct {
	dbConn *pgx.Conn
}

func NewProductService(conn *pgx.Conn) *Service {
	return &Service{dbConn: conn}
}

func (s Service) ListProducts(ctx context.Context) ([]product.Product, error) {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	ds := datastore.NewProductDataStore(tx)
	return ds.ListProducts(ctx)
}

func (s Service) GetProductById(id string, ctx context.Context) (product.Product, error) {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return product.Product{}, err
	}

	ds := datastore.NewProductDataStore(tx)
	return ds.GetProductById(id, ctx)
}