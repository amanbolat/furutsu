package cartsrv

import (
	"context"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/cart"
)

type Service struct {
	repo datastore.Repository
}

func NewCartService(repo datastore.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) CreateCart(userId string, ctx context.Context) (cart.Cart, error) {
	ds := datastore.NewCartDataStore(s.repo)
	return ds.CreateCart(userId, ctx)
}

// func (s Service) GetCart(userId string, ctx context.Context) (cart.Cart, error) {
// 	ds := datastore.NewCartDataStore(s.repo)
//
// 	c, err := ds.GetCartForUser(userId, ctx)
// 	if err != nil {
// 		return cart.Cart{}, err
// 	}
//
// 	discSrv := discountsrv.NewService(s.repo)
//
// 	discountedCart, err := discSrv.ApplyDiscounts(c, ctx)
// 	if err != nil {
// 		return cart.Cart{}, err
// 	}
//
// 	return discountedCart, nil
// }

// func (s Service) SetItemAmount(cartId string, productId string, ctx context.Context) (cart.Cart, error) {
//
//
// 	return c, nil
// }
//
// func (s Service) ApplyCoupon(userId string, coupon discount.Coupon, ctx context.Context) (cart.Cart, error) {
// 	var c cart.Cart
//
// 	return c, nil
// }