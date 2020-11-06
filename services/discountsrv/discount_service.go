package discountsrv

import (
	"context"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/jackc/pgx/v4"
)

type Service struct {
	dbConn *pgx.Conn
}

func NewService(conn *pgx.Conn) *Service {
	return &Service{dbConn: conn}
}

func (s Service) ApplyDiscounts(c cart.Cart, ctx context.Context) (cart.Cart, error) {
	ds := datastore.NewDiscountDataStore(s.dbConn)
	discounts, err := ds.ListDiscounts(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	for _, coupon := range c.Coupons {
		if coupon == nil {
			continue
		}

		if coupon.GetPercentage() == 0 {
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
		if discountSet == nil{
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