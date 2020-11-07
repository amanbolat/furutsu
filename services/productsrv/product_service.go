package productsrv

import (
	"context"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/product"
)

type Service struct {
	dbConn *pgx.Conn
}

func NewProductService(conn *pgx.Conn) *Service {
	return &Service{dbConn: conn}
}

func (s Service) ListProducts(ctx context.Context) ([]product.Product, error) {
	ds := datastore.NewProductDataStore(datastore.NewPgxConn(s.dbConn))
	return ds.ListProducts(ctx)
}

func (s Service) GetProductById(id string, ctx context.Context) (product.Product, error) {
	ds := datastore.NewProductDataStore(datastore.NewPgxConn(s.dbConn))
	return ds.GetProductById(id, ctx)
}
