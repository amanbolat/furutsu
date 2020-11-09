package discount

import (
	"github.com/amanbolat/furutsu/internal/cart"
)

type Discount struct {
	Name    string `json:"name"`
	Rule    Rule   `json:"rule"`
	Percent int    `json:"percent"`
}

func (d Discount) GetDiscountSetFor(items map[string]cart.Item) (*cart.ItemsSet, map[string]cart.Item) {
	discSet, leftItems := d.Rule.Check(items)
	if discSet != nil {
		return &cart.ItemsSet{
			Set:             discSet,
			DiscountPercent: d.Percent,
			DiscountName:    d.Name,
		}, leftItems
	}

	return nil, leftItems
}
