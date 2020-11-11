package paymentsrv

import (
	"context"
	"github.com/amanbolat/furutsu/internal/apperr"
	"github.com/amanbolat/furutsu/internal/payment"
	v "github.com/go-ozzo/ozzo-validation"

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
	payment.CardData
	OrderId string `json:"order_id"`
}

func (req PayForTheOrderRequest) Validate() error {
	return v.ValidateStruct(&req,
		v.Field(&req.OrderId, v.Required),
		v.Field(&req.CardData.Number, v.Required),
		v.Field(&req.CardData.CVC, v.Required),
		v.Field(&req.CardData.Holder, v.Required),
		v.Field(&req.CardData.Year, v.Required),
		v.Field(&req.CardData.Month, v.Required),
	)
}

func (s Service) PayForTheOrder(req PayForTheOrderRequest, ctx context.Context) error {
	err := req.Validate()
	if err != nil {
		return apperr.With(err, "your payment was rejected", "")
	}

	ordSrv := ordersrv.NewService(s.repo)
	return ordSrv.UpdateOrderStatus(req.OrderId, order.StatusPaid, ctx)
}
