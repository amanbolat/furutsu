package discountsrv

import (
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
)

type Service struct {

}

// func (s Service) ApplyDiscounts(c cart.Cart, discounts []discount.Discount) {
//
// }

func ApplyDiscountsToCart(c cart.Cart, discounts []discount.Discount) cart.Cart {
	leftItems := make(map[string]cart.ItemLine)
	for k, v := range c.Items {
		leftItems[k] = v
	}

	var discountSetArr []cart.DiscountSet

	for _, d := range discounts {
		discountSet, li := d.GetDiscountSetFor(leftItems)
		leftItems = li
		discountSetArr = append(discountSetArr, *discountSet)
	}

	copy(c.DiscountSets, discountSetArr)

	return c
}