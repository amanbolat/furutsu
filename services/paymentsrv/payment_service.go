package paymentsrv

import (
	"context"
	"github.com/amanbolat/furutsu/internal/payment"

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
	CardData payment.CardData `json:"card_data"`
	OrderId  string           `json:"order_id"`
}

func (s Service) PayForTheOrder(req PayForTheOrderRequest, ctx context.Context) error {
	ordSrv := ordersrv.NewService(s.repo)
	return ordSrv.UpdateOrderStatus(req.OrderId, order.StatusPaid, ctx)
}
