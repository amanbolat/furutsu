package discount

import (
	"github.com/amanbolat/furutsu/internal/cart"
)

type RuleType string

const (
	RuleTypeSet RuleType = "set"
	RuleTypeAll RuleType = "all"
)

type Rule interface {
	// Check checks if the given items satisfy the rule.
	// If true it returns set of product_ids and discount applied
	// also the set items left without discount
	Check(map[string]cart.ItemLine) (discountSet map[string]int, itemsLeft map[string]cart.ItemLine)
}

type RuleItemsAll struct {
	ProductID string
	Amount    int
}

func (r RuleItemsAll) Check(items map[string]cart.ItemLine) (discountSet map[string]int, itemsLeft map[string]cart.ItemLine) {
	for productId, item := range items {
		if r.ProductID == productId && item.Amount >= r.Amount {
			items[productId] = cart.ItemLine{
				Product: item.Product,
				Amount:  0,
			}
			return map[string]int{productId: item.Amount}, items
		}
	}

	return nil, items
}

type RuleItemsSet struct {
	ItemSet map[string]int
}

func (r RuleItemsSet) Check(items map[string]cart.ItemLine) (discountSet map[string]int, itemsLeft map[string]cart.ItemLine) {
	discSet := make(map[string]int)
	oldItems := make(map[string]cart.ItemLine)
	for k, v := range items {
		oldItems[k] = v
	}

recLoop:
	for {
		for productId, amountNeeded := range r.ItemSet {
			item, ok := items[productId]
			if !ok {
				break recLoop
			}

			if item.Amount < amountNeeded {
				break recLoop
			}

			discSet[productId] += amountNeeded
			items[productId] = cart.ItemLine{
				Product: item.Product,
				Amount:  item.Amount - amountNeeded,
			}
		}
	}

	if len(discSet) != len(r.ItemSet) {
		return nil, oldItems
	}

	return discSet, items
}

type Discount struct {
	Rule    Rule
	Percent int
}

func (d Discount) GetDiscountSetFor(items map[string]cart.ItemLine) (*cart.DiscountSet, map[string]cart.ItemLine) {
	discSet, leftItems := d.Rule.Check(items)
	if discSet != nil {
		return &cart.DiscountSet{
			ItemsSet: discSet,
			Discount: d.Percent,
		}, leftItems
	}

	return nil, leftItems
}
