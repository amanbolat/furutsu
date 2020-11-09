package ordersrv

import (
	"context"
	"github.com/amanbolat/furutsu/internal/apperr"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/order"
	"github.com/amanbolat/furutsu/services/cartsrv"
)

type Service struct {
	repo datastore.Repository
}

func NewService(repo datastore.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) ListOrders(userId string, ctx context.Context) ([]order.Order, error) {
	ds := datastore.NewOrderDataStore(s.repo)
	ol, err := ds.ListOrders(userId, ctx)
	if err != nil {
		return nil, err
	}

	return ol, nil
}

func (s Service) GetOrderById(id, userId string, ctx context.Context) (order.Order, error) {
	ds := datastore.NewOrderDataStore(s.repo)
	o, err := ds.GetOrderById(id, userId, ctx)
	if err != nil {
		return order.Order{}, err
	}

	return o, nil
}

// CreateOrder checkouts all items int the cart and creates a new order.
// It clears all the items in the cart and binds all the coupons to the order.
func (s Service) CreateOrder(userId string, ctx context.Context) (order.Order, error) {
	var newOrder order.Order

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return order.Order{}, err
	}

	cartSrv := cartsrv.NewCartService(tx)
	c, err := cartSrv.GetCart(userId, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return order.Order{}, err
	}

	// Check coupons
	for _, coupon := range c.Coupons {
		if coupon.IsUsed(c.Id) || coupon.IsExpired() {
			_ = tx.Rollback(ctx)
			return order.Order{}, apperr.New("some coupons are already invalid", "please checkout the cart again")
		}
	}

	cds := datastore.NewCartDataStore(tx)
	// Detach coupons from the cart
	for _, coupon := range c.Coupons {
		err = cds.DetachCouponFromCart(coupon.GetCode(), ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return order.Order{}, apperr.With(err, "failed to create the order", "")
		}
	}

	// Map items
	for _, item := range c.Items {
		oi := order.OrderItem{
			ProductName:        item.Product.Name,
			ProductDescription: item.Product.Description,
			Price:              item.Product.Price,
			Amount:             item.Amount,
		}
		newOrder.Items = append(newOrder.Items, oi)
	}
	newOrder.Total = c.Total()
	newOrder.Savings = c.TotalSavings()
	newOrder.TotalForPayment = c.TotalForPayment()
	newOrder.Status = order.StatusPending
	newOrder.UserId = userId

	// Clear cart
	err = cartSrv.ClearCart(c.Id, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return order.Order{}, err
	}

	ds := datastore.NewOrderDataStore(tx)
	resOrder, err := ds.CreateOrder(newOrder, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return order.Order{}, err
	}

	// Attach coupons
	for _, coupon := range c.Coupons {
		err = cds.AttachCouponToOrder(resOrder.Id, coupon.GetCode(), ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			return order.Order{}, apperr.With(err, "failed to create the order", "")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return order.Order{}, nil
	}

	return resOrder, nil
}

func (s Service) UpdateOrderStatus(orderId string, status order.Status, ctx context.Context) error {
	ds := datastore.NewOrderDataStore(s.repo)
	err := ds.UpdateOrderStatus(orderId, status, ctx)
	if err != nil {
		return err
	}

	return nil
}
