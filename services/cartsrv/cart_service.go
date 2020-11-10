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

// AddProductToCart adds a given amount of the product to cart as cart item
// or create new ones
func (s Service) AddProductToCart(productId, userId string, amount int, ctx context.Context) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}

	ds := datastore.NewCartDataStore(tx)
	cartId, err := ds.GetCartIdForUser(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	foundItem, err := ds.GetCartItem(cartId, productId, ctx)
	// If no cart item found we should create a new one,
	// so ErrNoRecords hasn't to be returned
	if err != nil && !errors.Is(err, datastore.ErrNoRecords) {
		_ = tx.Rollback(ctx)
		return err
	}

	newAmount := foundItem.Amount + amount

	if errors.Is(err, datastore.ErrNoRecords) {
		err = ds.CreateCartItem(cartId, productId, amount, ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	} else {
		err = ds.SetCartItemAmount(cartId, productId, newAmount, ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SetItemAmount sets a particular amount of cart items or removes them from the cart
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

	ds := datastore.NewCartDataStore(tx)
	c, err := s.GetCart(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, err
	}

	discountDs := datastore.NewDiscountDataStore(tx)
	foundCoupon, err := discountDs.GetCoupon(couponCode, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.With(err, "coupon code is not valid", "")
	}
	if foundCoupon.IsExpired() || foundCoupon.IsUsed() || !foundCoupon.IsAppliedToCart(c.Id) {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.New("coupon code is not valid", "")
	}

	// Check if there are items to which the coupon can be applied.
	// Coupons for the same product CANNOT be stacked together.
	dItems, _ := foundCoupon.Rule.Check(c.NonDiscountSetItems())
	var discAppliedAmount int
	for _, a := range dItems {
		discAppliedAmount += a
	}

	if discAppliedAmount < 1 {
		_ = tx.Rollback(ctx)
		return cart.Cart{}, apperr.New("no items in the cart to which the coupon could be applied", "")
	}

	err = ds.AttachCouponToCart(userId, couponCode, ctx)
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
	err = ds.DetachCouponFromCart(couponCode, ctx)
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
