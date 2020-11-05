package cartsrv

import (
	"context"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/jackc/pgx/v4"
)

type Service struct {
	dbConn *pgx.Conn
}

func NewCartService(conn *pgx.Conn) *Service {
	return &Service{dbConn: conn}
}

func (s Service) GetCart(userId string, ctx context.Context) (cart.Cart, error) {
	var c cart.Cart

	return c, nil
}

func (s Service) SetItemAmount(userId string, item cart.ItemLine, ctx context.Context) (cart.Cart, error) {
	var c cart.Cart

	return c, nil
}

func (s Service) ApplyCoupon(userId string, coupon discount.Coupon, ctx context.Context) (cart.Cart, error) {
	var c cart.Cart

	return c, nil
}