package datastore

import (
	"context"
	"time"

	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/georgysavva/scany/pgxscan"
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
	repo Repository
}

func NewDiscountDataStore(repo Repository) *DiscountDataStore {
	return &DiscountDataStore{repo: repo}
}

func (s DiscountDataStore) ListDiscounts(ctx context.Context) ([]discount.Discount, error) {
	var arr []DbDiscount
	err := pgxscan.Select(ctx, s.repo, &arr, `select * from discount`)
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
