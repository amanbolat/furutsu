package cart

import (
	"github.com/amanbolat/furutsu/internal/product"
)

type ItemLine struct {
	Product product.Product
	Amount int
}

type DiscountSet struct {
	// ItemsSet represent set of multiple items and their amounts
	// which share one discount
	ItemsSet map[string]int
	Discount int
}

type Cart struct {
	Items map[string]ItemLine
	DiscountSets []DiscountSet
}

func (c *Cart) AddProduct(p product.Product, amount int) {
	if c.Items == nil {
		c.Items = make(map[string]ItemLine)
	}

	il := ItemLine{
		Product: p,
		Amount:  amount,
	}

	if oldItem, ok := c.Items[p.ID]; ok {
		il.Amount += oldItem.Amount
		c.Items[p.ID] = il
	}

	c.Items[p.ID] = il
}

