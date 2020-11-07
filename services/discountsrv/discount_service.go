package discountsrv

import (
	"context"
	"time"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
)

type Service struct {
	repo datastore.Repository
}

func NewService(repo datastore.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) ApplyDiscounts(c cart.Cart, ctx context.Context) (cart.Cart, error) {
	ds := datastore.NewDiscountDataStore(s.repo)
	discounts, err := ds.ListDiscounts(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	for _, coupon := range c.Coupons {
		if coupon == nil {
			continue
		}

		// TODO: may be we should consider to put the check of expiration time
		// into the database because of the different timezones
		if coupon.GetPercentage() == 0 || coupon.GetExpireTime().Before(time.Now()) {
			continue
		}

		v, ok := coupon.(discount.Coupon)
		if !ok {
			continue
		}

		d := discount.Discount{
			Rule:    v.Rule,
			Percent: 0,
		}
		discounts = append(discounts, d)
	}

	newCart := ApplyDiscountsToCart(c, discounts)
	return newCart, nil
}

// ApplyDiscountsToCart applies all the discounts given to the cart items
// including coupon discounts and returns a new cart
func ApplyDiscountsToCart(c cart.Cart, discounts []discount.Discount) cart.Cart {
	leftItems := make(map[string]cart.Item)
	for k, v := range c.Items {
		leftItems[k] = v
	}

	var discountSetArr []cart.ItemsSet

	for _, d := range discounts {
		discountSet, li := d.GetDiscountSetFor(leftItems)
		if discountSet == nil {
			continue
		}
		leftItems = li
		discountSetArr = append(discountSetArr, *discountSet)
	}

	c.DiscountSets = make([]cart.ItemsSet, len(discountSetArr))
	copy(c.DiscountSets, discountSetArr)

	c.NonDiscountSet.Set = make(map[string]int)
	c.NonDiscountSet.Discount = 0
	for k, v := range leftItems {
		c.NonDiscountSet.Set[k] = v.Amount
	}

	return c
}
