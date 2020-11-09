package cartsrv_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"testing"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/services/cartsrv"
	"github.com/stretchr/testify/assert"
)

func TestService_GetCart(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
	assert.NoError(t, err)

	srv := cartsrv.NewCartService(datastore.NewPgxConn(conn))
	c, err := srv.GetCart("9a5f9080-978f-4260-8a45-1c3fbd30b6d4", context.Background())
	assert.NoError(t, err)

	for _, item := range c.Items {
		fmt.Println(item.Product.Name, item.Amount, "x", item.Product.Price)
	}

	for _, set := range c.DiscountSets {
		fmt.Printf("DiscountPercent: %v, [%d]\n", set.Set, set.DiscountPercent)
	}
	fmt.Println(c.TotalForPayment())
	fmt.Println(c.TotalSavings())
}

func TestService_SetItemAmount(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
	assert.NoError(t, err)

	srv := cartsrv.NewCartService(datastore.NewPgxConn(conn))
	c, err := srv.SetItemAmount("748cb518-1bd4-4f45-98ba-be016b39827e", "9a5f9080-978f-4260-8a45-1c3fbd30b6d4", 10, context.Background())
	assert.NoError(t, err)

	t.Log(c)
}
