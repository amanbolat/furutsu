package discount_test

import (
	"testing"

	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/stretchr/testify/assert"
)

func TestRuleItemsAll_Check(t *testing.T) {
	var testTable = []struct {
		name         string
		rule         discount.RuleItemsAll
		inItems      map[string]cart.Item
		outSet       map[string]int
		outLeftItems map[string]cart.Item
	}{
		{
			name: "7 apples discount - true",
			rule: discount.RuleItemsAll{
				ProductID: "apple",
				Amount:    7,
			},
			inItems: map[string]cart.Item{
				"apple": {Product: product.Product{
					ID: "apple",
				},
					Amount: 10,
				},
			},
			outSet: map[string]int{"apple": 10},
			outLeftItems: map[string]cart.Item{
				"apple": {Product: product.Product{
					ID: "apple",
				},
					Amount: 0,
				},
			},
		},
		{
			name: "7 apples discount - not enough amount",
			rule: discount.RuleItemsAll{
				ProductID: "apple",
				Amount:    7,
			},
			inItems: map[string]cart.Item{
				"apple": {Product: product.Product{
					ID: "apple",
				},
					Amount: 5,
				},
			},
			outSet: nil,
			outLeftItems: map[string]cart.Item{
				"apple": {Product: product.Product{
					ID: "apple",
				},
					Amount: 5,
				},
			},
		},
		{
			name: "7 apples discount - no apple item",
			rule: discount.RuleItemsAll{
				ProductID: "apple",
				Amount:    7,
			},
			inItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear",
				},
					Amount: 5,
				},
			},
			outSet: nil,
			outLeftItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear",
				},
					Amount: 5,
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			set, itemsLeft := testCase.rule.Check(testCase.inItems)
			assert.Equal(t, testCase.outSet, set)
			assert.Equal(t, testCase.outLeftItems, itemsLeft)
		})
	}
}

func TestRuleItemsSet_Check(t *testing.T) {
	bananaPearRule := discount.RuleItemsSet{
		ItemSet: map[string]int{
			"pear":   4,
			"banana": 2,
		},
	}

	var testTable = []struct {
		name         string
		rule         discount.RuleItemsSet
		inItems      map[string]cart.Item
		outSet       map[string]int
		outLeftItems map[string]cart.Item
	}{
		{
			name: "4 pears, 2 bananas - 1 set",
			rule: bananaPearRule,
			inItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 4,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 2,
				},
			},
			outSet: map[string]int{"pear": 4, "banana": 2},
			outLeftItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 0,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 0,
				},
			},
		},
		{
			name: "10 pears, 10 bananas - 2 sets",
			rule: bananaPearRule,
			inItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 10,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 10,
				},
			},
			outSet: map[string]int{"pear": 8, "banana": 4},
			outLeftItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 2,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 6,
				},
			},
		},
		{
			name: "7 pears, 1 bananas - 0 set",
			rule: bananaPearRule,
			inItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 7,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 1,
				},
			},
			outSet: nil,
			outLeftItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 7,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 1,
				},
			},
		},
		{
			name: "18 pears, 2 bananas discount - 1 set",
			rule: bananaPearRule,
			inItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 18,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 2,
				},
			},
			outSet: map[string]int{"pear": 4, "banana": 2},
			outLeftItems: map[string]cart.Item{
				"pear": {Product: product.Product{
					ID: "pear"},
					Amount: 14,
				},
				"banana": {Product: product.Product{
					ID: "banana"},
					Amount: 0,
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			set, leftItems := testCase.rule.Check(testCase.inItems)
			assert.Equal(t, testCase.outSet, set)
			assert.Equal(t, testCase.outLeftItems, leftItems)
		})
	}
}
