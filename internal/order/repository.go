package order

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *sqlcgen.Queries
}

var noOrderFound = errors.New("no orders found!")

func NewOrderRepositories(q *sqlcgen.Queries) *OrderRepository {
	return &OrderRepository{
		db: q,
	}
}

func (or *OrderRepository) CreateOrders(ctx context.Context, userId uuid.UUID, productId uuid.UUID, price float64) (model.Order, error) {
	newData, err := or.db.CreateOrder(ctx, sqlcgen.CreateOrderParams{
		UserID:       userId,
		ProductID:    productId,
		OrderedPrice: price,
	})
	if err != nil {
		return model.Order{}, err
	}
	return model.OrderDataToModel(newData), nil
}

func (or *OrderRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (model.Order, error) {
	data, err := or.db.UpdateOrderStatus(ctx, sqlcgen.UpdateOrderStatusParams{
		ID:     id,
		Status: status,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, noOrderFound
		}

		return model.Order{}, err
	}

	return model.OrderDataToModel(data), nil
}

func (or *OrderRepository) GetAllOrders(ctx context.Context, status *string) ([]model.Order, error) {
	data, err := or.db.ListOrders(ctx, status)
	if err != nil {
		return []model.Order{}, err
	}

	var listModel []model.Order = make([]model.Order, len(data))

	for i, v := range data {
		listModel[i] = model.MapGenToOrder(v.Order, v.User, v.Product)
	}

	return listModel, nil
}

func (or *OrderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {

	data, err := or.db.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, noOrderFound
		}

		return model.Order{}, err
	}

	return model.MapGenToOrder(data.Order, data.User, data.Product), nil
}

func (or *OrderRepository) GetOrderFromUsers(ctx context.Context, userId uuid.UUID, status *string) ([]model.Order, error) {
	data, err := or.db.ListOrdersByUser(ctx, sqlcgen.ListOrdersByUserParams{
		UserID: userId,
		Status: status,
	})

	if err != nil {
		return []model.Order{}, err
	}

	var listModel []model.Order = make([]model.Order, len(data))
	for i, v := range data {
		listModel[i] = model.MapGenToOrder(v.Order, v.User, v.Product)
	}

	return listModel, nil
}

func (or *OrderRepository) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	rows, err := or.db.SoftDeleteOrder(ctx, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return noOrderFound
	}
	return nil
}
