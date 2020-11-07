package discount

import (
	"encoding/json"
	"github.com/amanbolat/furutsu/internal/cart"
)

type Rule interface {
	// Check checks if the given items satisfy the rule.
	// If true it returns set of product_ids and discount applied
	// also the set items left without discount
	Check(map[string]cart.Item) (discountSet map[string]int, itemsLeft map[string]cart.Item)
	ToJSON() []byte
}

type RuleType string

const (
	RuleTypeSet RuleType = "set"
	RuleTypeAll RuleType = "all"
)

type RuleItemsAll struct {
	ProductID string
	Amount    int
}

func (r RuleItemsAll) ToJSON() []byte {
	m := make(map[string]int)
	m[r.ProductID] = r.Amount
	b, _ := json.Marshal(m)
	return b
}

func (r RuleItemsAll) Check(items map[string]cart.Item) (discountSet map[string]int, itemsLeft map[string]cart.Item) {
	for productId, item := range items {
		if r.ProductID == productId && item.Amount >= r.Amount {
			items[productId] = cart.Item{
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


func (r RuleItemsSet) ToJSON() []byte {
	b, _ := json.Marshal(r.ItemSet)
	return b
}

func (r RuleItemsSet) Check(items map[string]cart.Item) (discountSet map[string]int, itemsLeft map[string]cart.Item) {
	discSet := make(map[string]int)
	oldItems := make(map[string]cart.Item)
	for k, v := range items {
		oldItems[k] = v
	}

	var maxNeedProduct string
	var maxNeedAmount int
	for p, a := range r.ItemSet {
		if a > maxNeedAmount {
			maxNeedAmount = a
			maxNeedProduct = p
		}
	}

	item, ok := items[maxNeedProduct]
	if !ok {
		return nil, oldItems
	}

	maxSets := item.Amount / maxNeedAmount

	for productId, amountNeeded := range r.ItemSet {
		item, ok := items[productId]
		if !ok {
			break
		}

		var setsAdded int
		leftAmount := item.Amount
		for leftAmount >= amountNeeded && setsAdded < maxSets {
			leftAmount -= amountNeeded
			setsAdded += 1
			discSet[productId] += amountNeeded
			items[productId] = cart.Item{
				Product: item.Product,
				Amount:  leftAmount,
			}
		}
	}

	if len(discSet) != len(r.ItemSet) {
		return nil, oldItems
	}

	return discSet, items
}
