package datastore

import (
	"strconv"

	"github.com/amanbolat/furutsu/internal/discount"
)

func convertJsonbToRule(m map[string]interface{}) discount.Rule {
	var rule discount.Rule

	if m == nil || len(m) < 1 {
		return rule
	}

	if len(m) == 1 {
		r := discount.RuleItemsAll{}
		for k, v := range m {
			r.ProductID = k
			r.Amount = toInt(v)
		}
		rule = r
	} else {
		r := discount.RuleItemsSet{
			ItemSet: make(map[string]int),
		}
		for k, v := range m {
			r.ItemSet[k] = toInt(v)
		}
		rule = r
	}

	return rule
}

func toInt(data interface{}) int {
	switch v := data.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		return 0
	}
}
