package datastore_test

import (
	"context"
	"github.com/jackc/pgx/v4"
	"testing"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/stretchr/testify/assert"
)

func TestProductDataStore_ListProducts(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
	assert.NoError(t, err)

	ds := datastore.NewProductDataStore(datastore.NewPgxConn(conn))

	products, err := ds.ListProducts(context.Background())
	assert.NoError(t, err)
	t.Log(products)
}
