package ordersrv

import (
	"context"

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

func (s Service) ListOrders() {

}

func (s Service) GetOrder() {

}

// CreateOrder checkouts all items int the cart and creates a new order
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

	err = cartSrv.ClearCart(c.Id, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return order.Order{}, err
	}

	ds := datastore.NewOrderDataStore(tx)
	resOrder, err := ds.CreateOrder(newOrder, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return order.Order{}, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return order.Order{}, nil
	}

	return resOrder, nil
}

func (s Service) UpdateOrderStatus(orderId string, status order.Status, ctx context.Context) error {
	ds := datastore.NewOrderDataStore(s.repo)
	err := ds.UpdateStatus(orderId, status, ctx)
	if err != nil {
		return err
	}

	return nil
}
