package datastore_test

import (
	"context"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductDataStore_ListProducts(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
	assert.NoError(t, err)

	tx, err := conn.Begin(context.Background())
	defer tx.Commit(context.Background())
	assert.NoError(t, err)
	ds := datastore.NewProductDataStore(tx)

	products, err := ds.ListProducts(context.Background())
	assert.NoError(t, err)
	t.Log(products)
}
