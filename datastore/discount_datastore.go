package datastore

import (
	"context"
	"database/sql"
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

func (d DbDiscount) ToDiscount() discount.Discount {
	rule := convertJsonbToRule(d.Rule)
	return discount.Discount{
		Name:    d.Name,
		Rule:    rule,
		Percent: d.Percent,
	}
}

type DbCoupon struct {
	ID              string
	UserId          sql.NullString
	Code            string
	Name            string
	CartId          sql.NullString
	Rule            map[string]interface{}
	DiscountPercent int
	ExpireAt        time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c DbCoupon) ToCoupon() discount.Coupon {
	rule := convertJsonbToRule(c.Rule)

	return discount.Coupon{
		ID:      c.ID,
		Code:    c.Code,
		Name:    c.Name,
		Rule:    rule,
		Percent: c.DiscountPercent,
		Expire:  c.ExpireAt,
	}
}

type DiscountDataStore struct {
	repo Repository
}

func NewDiscountDataStore(repo Repository) *DiscountDataStore {
	return &DiscountDataStore{repo: repo}
}

func (s DiscountDataStore) CreateDiscount(d discount.Discount, ctx context.Context) (discount.Discount, error) {
	rows, err := s.repo.Query(ctx,
		`INSERT INTO discount (name, rule, percent)
VALUES ($1, $2, $3) RETURNING *`,
		d.Name, d.Rule.ToJSON(), d.Percent,
	)
	if err != nil {
		return discount.Discount{}, err
	}
	defer rows.Close()

	var dbDiscount DbDiscount
	err = pgxscan.ScanOne(&dbDiscount, rows)
	if err != nil {
		return discount.Discount{}, err
	}

	return dbDiscount.ToDiscount(), nil
}

func (s DiscountDataStore) ListDiscounts(ctx context.Context) ([]discount.Discount, error) {
	var arr []DbDiscount
	err := pgxscan.Select(ctx, s.repo, &arr, `select * from discount`)
	if err != nil {
		return nil, err
	}

	discounts := make([]discount.Discount, len(arr))
	for i, d := range arr {
		discounts[i] = d.ToDiscount()
	}

	return discounts, nil
}

func (s DiscountDataStore) CreateCoupon(c discount.Coupon, ctx context.Context) (discount.Coupon, error) {
	rows, err := s.repo.Query(ctx,
		`INSERT INTO coupon (code, name, expire_at, rule, discount_percent) 
VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		c.Code, c.Name, c.Expire, c.Rule.ToJSON(), c.Percent,
	)
	if err != nil {
		return discount.Coupon{}, err
	}
	defer rows.Close()

	var dbCoupon DbCoupon
	err = pgxscan.ScanOne(&dbCoupon, rows)
	if err != nil {
		return discount.Coupon{}, err
	}

	return dbCoupon.ToCoupon(), nil
}
