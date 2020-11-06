package datastore

import "github.com/amanbolat/furutsu/internal/discount"

func convertJsonbToRule(m map[string]interface{}) discount.Rule {
	var rule discount.Rule

	if m == nil || len(m) < 1 {
		return rule
	}

	if len(m) == 1 {
		r := discount.RuleItemsAll{}
		for k, v := range m {
			r.ProductID = k
			r.Amount = v.(int)
		}
		rule = r
	} else {
		r := discount.RuleItemsSet{
			ItemSet: make(map[string]int),
		}
		for k, v := range m {
			r.ItemSet[k] = v.(int)
		}
		rule = r
	}

	return rule
}
