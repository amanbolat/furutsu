package discountsrv_test

import (
	"testing"

	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/amanbolat/furutsu/services/discountsrv"
	"github.com/stretchr/testify/assert"
)

func TestApplyDiscountsToCart(t *testing.T) {
	c1 := cart.Cart{
		Items: map[string]cart.Item{
			"pear": {
				Product: product.Product{
					ID:    "pear",
					Price: 100,
				},
				Amount: 10,
			},
			"banana": {
				Product: product.Product{
					ID:    "banana",
					Price: 100,
				},
				Amount: 10,
			},
			"apple": {
				Product: product.Product{
					ID:    "apple",
					Price: 100,
				},
				Amount: 10,
			},
			"orange": {
				Product: product.Product{
					ID:    "orange",
					Price: 100,
				},
				Amount: 10,
			},
		},
	}

	d1 := []discount.Discount{
		{
			Rule: discount.RuleItemsAll{
				ProductID: "apple",
				Amount:    7,
			},
			Percent: 10,
		},
		{
			Rule: discount.RuleItemsSet{
				ItemSet: map[string]int{
					"pear":   4,
					"banana": 2,
				},
			},
			Percent: 30,
		},
	}

	resCart := discountsrv.ApplyDiscountsToCart(c1, d1)
	expectedDiscSet := []cart.ItemsSet{
		{
			Set:             map[string]int{"apple": 10},
			DiscountPercent: 10,
		},
		{
			Set:             map[string]int{"pear": 8, "banana": 4},
			DiscountPercent: 30,
		},
	}
	assert.Equal(t, c1.Items, resCart.Items)
	assert.Equal(t, expectedDiscSet, resCart.DiscountSets)
	assert.Equal(t, 3540, resCart.TotalForPayment())
	assert.Equal(t, 460, resCart.TotalSavings())
}
