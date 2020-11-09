package cartsrv

import (
	"context"
	"errors"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/apperr"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/services/discountsrv"
)

type Service struct {
	repo datastore.Repository
}

func NewCartService(repo datastore.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) WithTx(tx datastore.RepoTx) *Service {
	return NewCartService(tx)
}

func (s Service) CreateCart(userId string, ctx context.Context) (cart.Cart, error) {
	ds := datastore.NewCartDataStore(s.repo)
	return ds.CreateCart(userId, ctx)
}

func (s Service) GetCart(userId string, ctx context.Context) (cart.Cart, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	ds := datastore.NewCartDataStore(tx)

	c, err := ds.GetCartForUser(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	discSrv := discountsrv.NewService(tx)
	discountedCart, err := discSrv.ApplyDiscounts(c, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	return discountedCart, nil
}

// SetItemAmount used to add, remove items from the cart
// or to change its amount
func (s Service) SetItemAmount(productId, userId string, amount int, ctx context.Context) (cart.Cart, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	ds := datastore.NewCartDataStore(tx)
	cartId, err := ds.GetCartIdForUser(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	_, err = ds.GetCartItem(cartId, productId, ctx)
	// If no cart item found we should create a new one,
	// so ErrNoRecords hasn't to be returned
	if err != nil && !errors.Is(err, datastore.ErrNoRecords) {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	// Create item if there isn't any
	if errors.Is(err, datastore.ErrNoRecords) {
		err = ds.CreateCartItem(cartId, productId, amount, ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return cart.Cart{}, err
		}
		// Delete item if amount argument is less than 1
	} else if amount < 1 {
		err = ds.DeleteCartItem(cartId, productId, ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return cart.Cart{}, err
		}
		// Set new amount
	} else {
		err = ds.SetCartItemAmount(cartId, productId, amount, ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return cart.Cart{}, err
		}
	}

	// Call GetCart to apply all discounts and coupons
	c, err := s.WithTx(tx).GetCart(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	return c, nil
}

func (s Service) ApplyCoupon(userId, couponCode string, ctx context.Context) (cart.Cart, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	discountDs := datastore.NewDiscountDataStore(tx)
	foundCoupon, err := discountDs.GetCoupon(couponCode, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.With(err, "couldn't find the coupon", "")
	}
	if foundCoupon.IsExpired() || foundCoupon.IsUsed() {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.New("coupon is not valid anymore", "")
	}

	ds := datastore.NewCartDataStore(tx)
	c, err := ds.GetCartForUser(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	dItems, _ := foundCoupon.Rule.Check(c.Items)
	if dItems == nil || len(dItems) < 1 {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.New("no items in the cart to which the coupon could be applied", "")
	}

	err = ds.AttachCoupon(userId, couponCode, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, nil
	}

	c, err = s.WithTx(tx).GetCart(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return cart.Cart{}, nil
	}

	return c, nil
}

func (s Service) DetachCoupon(userId, couponCode string, ctx context.Context) (cart.Cart, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return cart.Cart{}, err
	}

	ds := datastore.NewCartDataStore(tx)
	err = ds.DetachCoupon(couponCode, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, nil
	}

	c, err := s.WithTx(tx).GetCart(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return cart.Cart{}, nil
	}

	return c, nil
}

func (s Service) ClearCart(cartId string, ctx context.Context) error {
	ds := datastore.NewCartDataStore(s.repo)
	err := ds.ClearCart(cartId, ctx)
	if err != nil {
		return err
	}

	return nil
}
