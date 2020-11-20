package discount_test

import (
	"testing"

	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/stretchr/testify/assert"
)

func TestCoupon_IsApplicableToItems(t *testing.T) {
	cp := discount.Coupon{
		Rule: discount.RuleItemsAll{
			ProductID: "orange",
			Amount:    0,
		},
		Percent: 10,
	}

	var orangeCart = cart.Cart{
		Items: map[string]cart.Item{
			"orange": {
				Product: product.Product{
					ID:   "orange",
					Name: "orange",
				},
				Amount: 11,
			},
		},
	}

	var appleCart = cart.Cart{
		Items: map[string]cart.Item{
			"apple": {
				Product: product.Product{
					ID:   "apple",
					Name: "apple",
				},
				Amount: 11,
			},
		},
	}

	assert.Equal(t, true, cp.IsApplicableToItems(orangeCart.Items))
	assert.Equal(t, false, cp.IsApplicableToItems(appleCart.Items))
}
