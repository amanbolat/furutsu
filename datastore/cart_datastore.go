package datastore

import (
	"context"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/georgysavva/scany/pgxscan"
	"time"
)

type CartDataStore struct {
	querier pgxscan.Querier
}

type DbCartItem struct {
	Id string
	CartId string
	ProductId string
	ProductName string
	ProductPrice int
	ProductDescription string
	CreatedAt time.Time
	UpdatedAt time.Time
	Amount int
}

type DbCart struct {
	Id string
	UserId string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (dc DbCart) ToCart() cart.Cart {
	return cart.Cart{
		Id:             dc.Id,
	}
}

type DbCoupon struct {
	ID      string
	Code    string
	Name    string
	CartId string
	Rule    map[string]interface{}
	DiscountPercent int
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCartDataStore(q pgxscan.Querier) *CartDataStore {
	return &CartDataStore{querier: q}
}

func (s CartDataStore) CreateCart(userId string, ctx context.Context) (cart.Cart, error) {
	rows, err := s.querier.Query(ctx, `INSERT INTO cart (user_id) VALUES ($1)`, userId)
	if err != nil {
		return cart.Cart{}, err
	}
	defer rows.Close()

	var dc DbCart
	err = pgxscan.ScanRow(&dc, rows)
	if err != nil {
		return cart.Cart{}, err
	}

	return dc.ToCart(), nil
}

func (s CartDataStore) GetCartForUser(userid string, ctx context.Context) (cart.Cart, error) {
	var dbCart DbCart
	var dbCartItems []DbCartItem
	var coupons []DbCoupon

	err := pgxscan.Get(ctx, s.querier, &dbCart, `SELECT * FROM cart WHERE user_id = $1`, userid)
	if err != nil {
		return cart.Cart{}, err
	}

	err = pgxscan.Select(ctx, s.querier, &dbCartItems, `
SELECT
       cart_item.id AS id,
       product_id,
       amount,
       p.name AS product_name,
       p.price AS product_price,
       p.description AS product_description,
       cart_item.created_at as created_at,
       cart_item.updated_at as updated_at
FROM cart_item
JOIN cart c ON c.id = cart_item.cart_id
JOIN product p ON p.id = cart_item.product_id
WHERE c.user_id = $1
`, userid)
	if err != nil {
		return cart.Cart{}, err
	}

	err = pgxscan.Select(ctx, s.querier, &coupons, `SELECT * FROM coupon WHERE cart_id = `, dbCart.Id)

	c := cart.Cart{
		Id:             dbCart.Id,
		Items: make(map[string]cart.Item, len(dbCartItems)),
		Coupons: make([]cart.Coupon, len(coupons)),
	}

	for _, item := range dbCartItems {
		c.Items[item.ProductId] = cart.Item{
			Id:      item.Id,
			Product: product.Product{
				ID:          item.ProductId,
				Name:        item.ProductName,
				Price:       item.ProductPrice,
				Description: item.ProductDescription,
			},
			Amount:  item.Amount,
		}
	}

	for i, coupon := range coupons {
		rule := convertJsonbToRule(coupon.Rule)

		c.Coupons[i] = discount.Coupon{
			ID:      coupon.ID,
			Code:    coupon.Code,
			Name:    coupon.Name,
			Rule:    rule,
			Percent: coupon.DiscountPercent,
			Expire:  coupon.ExpireAt,
		}
	}

	return c, nil
}

// func (s CartDataStore) UpdateCartItem(cartId string, productId string, ctx context.Context) (cart.Cart, error) {
//
//
// 	return newCart, nil
// }