package order

import (
	"context"
	"encoding/json"
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
		listModel[i] = MapListOrdersToOrder(v)
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

	return mapGetOrderByIDToOrder(data), nil
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
		listModel[i] = MapListOrdersByUserToOrder(v)
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

func parseOrderPayments(data []byte) []model.Payment {
	var summaries []orderPaymentSummary

	if err := json.Unmarshal(data, &summaries); err != nil {
		return nil
	}

	payments := make([]model.Payment, 0, len(summaries))
	for _, s := range summaries {
		payments = append(payments, model.Payment{
			ID:     s.ID,
			Method: s.Method,
			Status: s.Status,
			Amount: s.Amount,
		})
	}

	return payments
}

func MapListOrdersToOrder(row sqlcgen.ListOrdersRow) model.Order {
	user := model.MapToUserModel(row.User)
	product := model.MapToProductModel(row.Product)

	payments := parseOrderPayments(row.Payments)

	return model.Order{
		ID:           row.Order.ID,
		UserID:       row.Order.UserID,
		ProductID:    row.Order.ProductID,
		OrderedPrice: row.Order.OrderedPrice,
		Status:       row.Order.Status,
		CreatedAt:    row.Order.CreatedAt,
		UpdatedAt:    row.Order.UpdatedAt,
		DeletedAt:    row.Order.DeletedAt,
		User:         &user,
		Product:      &product,
		Payment:      payments,
	}
}

func MapListOrdersByUserToOrder(row sqlcgen.ListOrdersByUserRow) model.Order {
	user := model.MapToUserModel(row.User)
	product := model.MapToProductModel(row.Product)

	payments := parseOrderPayments(row.Payments)

	return model.Order{
		ID:           row.Order.ID,
		UserID:       row.Order.UserID,
		ProductID:    row.Order.ProductID,
		OrderedPrice: row.Order.OrderedPrice,
		Status:       row.Order.Status,
		CreatedAt:    row.Order.CreatedAt,
		UpdatedAt:    row.Order.UpdatedAt,
		DeletedAt:    row.Order.DeletedAt,
		User:         &user,
		Product:      &product,
		Payment:      payments,
	}
}

func mapGetOrderByIDToOrder(row sqlcgen.GetOrderByIDRow) model.Order {
	user := model.MapToUserModel(row.User)
	product := model.MapToProductModel(row.Product)

	var payments []model.Payment
	err := json.Unmarshal(row.Payments, &payments)
	if err != nil {
		panic(err)
	}

	return model.Order{
		ID:           row.Order.ID,
		UserID:       row.Order.UserID,
		ProductID:    row.Order.ProductID,
		OrderedPrice: row.Order.OrderedPrice,
		Status:       row.Order.Status,
		CreatedAt:    row.Order.CreatedAt,
		UpdatedAt:    row.Order.UpdatedAt,
		DeletedAt:    row.Order.DeletedAt,
		User:         &user,
		Product:      &product,
		Payment:      payments,
	}
}
