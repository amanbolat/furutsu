package datastore

import (
	"context"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/georgysavva/scany/pgxscan"
	"time"
)

type DbDiscount struct {
	Id        string
	Name      string
	Rule      map[string]interface{}
	Percent   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DiscountDataStore struct {
	querier pgxscan.Querier
}

func NewDiscountDataStore(q pgxscan.Querier) *DiscountDataStore {
	return &DiscountDataStore{querier: q}
}

func (s DiscountDataStore) ListDiscounts(ctx context.Context) ([]discount.Discount, error) {
	var arr []DbDiscount
	err := pgxscan.Select(ctx, s.querier, &arr, `select * from discount`)
	if err != nil {
		return nil, err
	}

	discounts := make([]discount.Discount, len(arr))
	for i, d := range arr {
		rule := convertJsonbToRule(d.Rule)
		discounts[i] = discount.Discount{
			Rule:    rule,
			Percent: d.Percent,
		}
	}

	return discounts, nil
}
