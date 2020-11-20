package cart_test

import (
	"testing"

	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCart = cart.Cart{
	Id:     "01",
	UserId: "02",
	Items: map[string]cart.Item{
		"apple": {
			Id: "apple",
			Product: product.Product{
				ID:    "apple",
				Name:  "apple",
				Price: 100,
			},
			Amount: 10,
		},
		"pear": {
			Id: "pear",
			Product: product.Product{
				ID:    "pear",
				Name:  "pear",
				Price: 150,
			},
			Amount: 30,
		},
	},
	DiscountSets:   nil,
	NonDiscountSet: cart.ItemsSet{},
	Coupons:        nil,
}

func TestCart_MarshalJSON(t *testing.T) {
	b, err := testCart.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestCart_Total(t *testing.T) {
	assert.Equal(t, 5500, testCart.Total())
}

func TestCart_TotalSavings(t *testing.T) {
	assert.Equal(t, 0, testCart.TotalSavings())
}

func TestCart_TotalForPayment(t *testing.T) {
	assert.Equal(t, 0, testCart.TotalSavings())
}

func TestCart_SetProductAmount(t *testing.T) {

}
