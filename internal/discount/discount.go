package discount

import (
	"github.com/amanbolat/furutsu/internal/cart"
)

type Discount struct {
	Name string
	Rule    Rule
	Percent int
}

func (d Discount) GetDiscountSetFor(items map[string]cart.Item) (*cart.ItemsSet, map[string]cart.Item) {
	discSet, leftItems := d.Rule.Check(items)
	if discSet != nil {
		return &cart.ItemsSet{
			Set:      discSet,
			Discount: d.Percent,
		}, leftItems
	}

	return nil, leftItems
}
