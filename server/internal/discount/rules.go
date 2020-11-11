package discount

import (
	"encoding/json"
	"github.com/amanbolat/furutsu/internal/cart"
	"sort"
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
	ProductID string `json:"product_id"`
	Amount    int    `json:"amount"`
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
	ItemSet map[string]int `json:"item_set"`
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

	var setCount []int
	for pId, need := range r.ItemSet {
		i := items[pId]
		maxSets := i.Amount / need
		setCount = append(setCount, maxSets)
	}
	sort.Ints(setCount)

	if len(setCount) < 1 {
		return nil, oldItems
	}

	minSet := setCount[0]
	if len(setCount) != len(r.ItemSet) || minSet < 1 {
		return nil, oldItems
	}

	for pId, need := range r.ItemSet {
		discSet[pId] += need * minSet
		leftAmount := items[pId].Amount - discSet[pId]
		items[pId] = cart.Item{
			Product: items[pId].Product,
			Amount:  leftAmount,
		}
	}

	if len(discSet) != len(r.ItemSet) {
		return nil, oldItems
	}

	return discSet, items
}
