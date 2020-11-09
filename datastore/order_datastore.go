package datastore

import (
	"context"
	"database/sql"
	"time"

	"github.com/amanbolat/furutsu/internal/order"
	"github.com/georgysavva/scany/pgxscan"
)

type OrderDataStore struct {
	querier pgxscan.Querier
}

func NewOrderDataStore(q pgxscan.Querier) *OrderDataStore {
	return &OrderDataStore{querier: q}
}

type DbOrderItem struct {
	Id                 string
	OrderId            string
	ProductName        string
	ProductDescription sql.NullString
	Price              int
	Amount             int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (i DbOrderItem) ToOrderItem() order.OrderItem {
	return order.OrderItem{
		Id:                 i.Id,
		OrderId:            i.OrderId,
		ProductName:        i.ProductName,
		ProductDescription: i.ProductDescription.String,
		Price:              i.Price,
		Amount:             i.Amount,
		CreatedAt:          i.CreatedAt,
		UpdatedAt:          i.UpdatedAt,
	}
}

type DbOrder struct {
	Id              string
	UserId          string
	Status          string
	Total           int
	Savings         int
	TotalForPayment int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (o DbOrder) ToOrder() order.Order {
	return order.Order{
		Id:              o.Id,
		UserId:          o.UserId,
		Status:          order.Status(o.Status),
		Total:           o.Total,
		Savings:         o.Savings,
		TotalForPayment: o.TotalForPayment,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

func (d OrderDataStore) GetOrderById(id, userId string, ctx context.Context) (order.Order, error) {
	var dbOrder DbOrder
	err := pgxscan.Get(ctx, d.querier, &dbOrder, `SELECT * FROM "order" WHERE id = $1 AND user_id = $2`, id, userId)
	if err != nil {
		return order.Order{}, err
	}

	var dbOrderItems []DbOrderItem
	err = pgxscan.Select(ctx, d.querier, &dbOrderItems, `SELECT * FROM order_item WHERE order_id = $1`, id)
	if err != nil {
		return order.Order{}, err
	}

	resOrder := dbOrder.ToOrder()

	for _, i := range dbOrderItems {
		resOrder.Items = append(resOrder.Items, i.ToOrderItem())
	}

	return resOrder, nil
}

func (d OrderDataStore) ListOrders(userId string, ctx context.Context) ([]order.Order, error) {
	var dbOrders []DbOrder
	err := pgxscan.Select(ctx, d.querier, &dbOrders, `SELECT * FROM "order" WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	var resOrders []order.Order
	for _, o := range dbOrders {
		resOrders = append(resOrders, o.ToOrder())
	}

	return resOrders, nil
}

func (d OrderDataStore) CreateOrder(ord order.Order, ctx context.Context) (order.Order, error) {
	var resOrder order.Order

	rows, err := d.querier.Query(ctx,
		`INSERT INTO "order" (user_id, status, total, savings, total_for_payment)
VALUES ($1, $2, $3, $4, $5) RETURNING *`, ord.UserId, ord.Status, ord.Total, ord.Savings, ord.TotalForPayment)
	if err != nil {
		return order.Order{}, err
	}
	defer rows.Close()

	var dbOrder DbOrder
	err = pgxscan.ScanOne(&dbOrder, rows)
	if err != nil {
		return order.Order{}, err
	}
	resOrder = dbOrder.ToOrder()

	for _, item := range ord.Items {
		r, err := d.querier.Query(ctx,
			`INSERT INTO order_item (product_name, product_description, order_id, price, amount)
VALUES ($1, $2, $3, $4, $5) RETURNING *`, item.ProductName, item.ProductDescription, dbOrder.Id, item.Price, item.Amount)
		if err != nil {
			return order.Order{}, err
		}
		var dbOrderItem DbOrderItem
		err = pgxscan.ScanOne(&dbOrderItem, r)
		if err != nil {
			return order.Order{}, err
		}
		resOrder.Items = append(resOrder.Items, dbOrderItem.ToOrderItem())

		r.Close()
	}

	return resOrder, nil
}

func (d OrderDataStore) UpdateOrderStatus(orderId string, status order.Status, ctx context.Context) error {
	rows, err := d.querier.Query(ctx, `UPDATE "order" SET status = $1 WHERE id = $2 AND status <> $1`, status, orderId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNoRecords
	}

	return nil
}
