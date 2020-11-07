package paymentsrv

import (
	"context"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/order"
	"github.com/amanbolat/furutsu/services/ordersrv"
)

type Service struct {
	repo datastore.Repository
}

func NewService(repo datastore.Repository) *Service {
	return &Service{repo: repo}
}

type PayForTheOrderRequest struct {
	CardNumber   int
	HolderName   string
	CardExpireAt string
	CVC          int
	OrderId      string
}

func (s Service) PayForTheOrder(req PayForTheOrderRequest, ctx context.Context) error {
	ordSrv := ordersrv.NewService(s.repo)
	return ordSrv.UpdateOrderStatus(req.OrderId, order.StatusPaid, ctx)
}
