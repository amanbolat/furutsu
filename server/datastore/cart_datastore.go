package datastore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type CartDataStore struct {
	querier pgxscan.Querier
}

type DbCartItem struct {
	Id                 string
	CartId             string
	ProductId          string
	ProductName        string
	ProductPrice       int
	ProductDescription sql.NullString
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Amount             int
}

func (i DbCartItem) ToCartItem() cart.Item {
	return cart.Item{
		Id: i.Id,
		Product: product.Product{
			ID:          i.ProductId,
			Name:        i.ProductName,
			Price:       i.ProductPrice,
			Description: i.ProductDescription.String,
		},
		Amount: i.Amount,
	}
}

type DbCart struct {
	Id        string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (dc DbCart) ToCart() cart.Cart {
	return cart.Cart{
		Id:     dc.Id,
		UserId: dc.UserId,
	}
}

func NewCartDataStore(q pgxscan.Querier) *CartDataStore {
	return &CartDataStore{querier: q}
}

func (s CartDataStore) CreateCart(userId string, ctx context.Context) (cart.Cart, error) {
	rows, err := s.querier.Query(ctx, `INSERT INTO cart (user_id) VALUES ($1) RETURNING *`, userId)
	if err != nil {
		return cart.Cart{}, err
	}
	defer rows.Close()

	var dc DbCart
	err = pgxscan.ScanOne(&dc, rows)
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
WHERE c.user_id = $1;`, userid)
	if err != nil {
		return cart.Cart{}, err
	}

	err = pgxscan.Select(ctx, s.querier, &coupons, `SELECT * FROM coupon WHERE cart_id = $1`, dbCart.Id)
	if err != nil {
		return cart.Cart{}, err
	}

	c := dbCart.ToCart()
	c.Items = make(map[string]cart.Item, len(dbCartItems))
	c.Coupons = make([]cart.Coupon, len(coupons))

	for _, item := range dbCartItems {
		c.Items[item.ProductId] = item.ToCartItem()
	}

	for i, coupon := range coupons {
		c.Coupons[i] = coupon.ToCoupon()
	}

	return c, nil
}

func (s CartDataStore) CreateCartItem(cartId, productId string, amount int, ctx context.Context) error {
	rows, err := s.querier.Query(ctx,
		`INSERT INTO cart_item (cart_id, product_id, amount) VALUES ($1, $2, $3)`,
		cartId, productId, amount)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s CartDataStore) GetCartItem(cartId, productId string, ctx context.Context) (cart.Item, error) {
	var item DbCartItem

	err := pgxscan.Get(ctx, s.querier, &item,
		`SELECT * FROM cart_item WHERE cart_id = $1 AND product_id = $2`,
		cartId,
		productId)
	if errors.Is(err, pgx.ErrNoRows) {
		return cart.Item{}, ErrNoRecords
	}

	if err != nil {
		return cart.Item{}, err
	}

	return item.ToCartItem(), nil
}

func (s CartDataStore) GetCartIdForUser(userId string, ctx context.Context) (string, error) {
	var id string
	err := pgxscan.Get(ctx, s.querier, &id, `SELECT id FROM cart WHERE user_id = $1`, userId)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s CartDataStore) DeleteCartItem(cartId, productId string, ctx context.Context) error {
	rows, err := s.querier.Query(ctx,
		`DELETE FROM cart_item WHERE cart_id = $1 AND product_id = $2`,
		cartId,
		productId)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s CartDataStore) SetCartItemAmount(cartId, productId string, amount int, ctx context.Context) error {
	rows, err := s.querier.Query(ctx,
		`UPDATE cart_item SET amount = $1 WHERE cart_id = $2 AND product_id = $3 RETURNING id`,
		amount,
		cartId,
		productId)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s CartDataStore) AttachCouponToCart(userId, couponCode string, ctx context.Context) error {
	rows, err := s.querier.Query(ctx,
		`
WITH c AS (
    SELECT id
    FROM cart
    WHERE user_id = $1
)
UPDATE coupon
SET cart_id = c.id
FROM c
WHERE coupon.code = $2
  AND coupon.expire_at > current_timestamp`,
		userId,
		couponCode,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s CartDataStore) DetachCouponFromCart(couponCode string, ctx context.Context) error {
	rows, err := s.querier.Query(ctx, `
UPDATE coupon
SET cart_id = NULL
WHERE code = $1
`, couponCode)

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}

func (s CartDataStore) ClearCart(cartId string, ctx context.Context) error {
	rows, err := s.querier.Query(ctx, `DELETE FROM cart_item WHERE cart_id = $1`, cartId)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (s CartDataStore) AttachCouponToOrder(orderId, couponCode string, ctx context.Context) error {
	rows, err := s.querier.Query(ctx,
		`
UPDATE coupon
SET order_id = $1
WHERE code = $2`,
		orderId,
		couponCode,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
