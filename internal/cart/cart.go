package cart

import (
	"github.com/amanbolat/furutsu/internal/product"
)

type Item struct {
	Id string
	Product product.Product
	Amount int
}

// ItemSet is set of items in cart which could have discounts
// applied or not. Used for calculating totals and convenient
// representation on client side
type ItemsSet struct {
	// Set represent set of multiple items and their amounts
	// which share one discount
	Set      map[string]int
	Discount int
}

type Coupon interface {
	GetPercentage() int
	GetName() string
}

type Cart struct {
	Id string
	// Items is map items as of product_id:CartItem
	Items map[string]Item
	DiscountSets []ItemsSet
	NonDiscountSet ItemsSet
	Coupons []Coupon
}

func (c Cart) TotalSavings() int {
	var total int
	for _, set := range c.DiscountSets {
		for productId, amount := range set.Set {
			price := c.Items[productId].Product.Price
			toPay := price * amount
			saved := toPay * set.Discount / 100
			total += saved
		}
	}

	return total
}

func (c Cart) TotalForPayment() int {
	var total int

	for _, item := range c.Items {
		toPay := item.Product.Price * item.Amount
		total += toPay
	}

	total -= c.TotalSavings()
	return total
}

func (c *Cart) SetProductAmount(p product.Product, amount int) {
	if c.Items == nil {
		c.Items = make(map[string]Item)
	}

	if amount < 1 {
		delete(c.Items, p.ID)
	}

	il := Item{
		Product: p,
		Amount:  amount,
	}

	if oldItem, ok := c.Items[p.ID]; ok {
		il.Amount += oldItem.Amount
		c.Items[p.ID] = il
	}

	c.Items[p.ID] = il
}

