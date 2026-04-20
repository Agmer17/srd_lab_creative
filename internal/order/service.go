package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type OrderService struct {
	repo           *OrderRepository
	productService *product.ProductService
}

func NewOrderService(rp *OrderRepository, psvc *product.ProductService) *OrderService {
	return &OrderService{
		repo:           rp,
		productService: psvc,
	}
}

func (osv *OrderService) GetAllOrders(ctx context.Context, status string) ([]OrderListDTO, *shared.ErrorResponse) {
	var tempStatus *string = nil
	if status != "" {
		tempStatus = &status
	}

	data, err := osv.repo.GetAllOrders(ctx, tempStatus)
	if err != nil {
		return []OrderListDTO{}, shared.NewErrorResponse(500, "something wrong while getting order data, try another time")
	}

	responseData := orderListModelToDto(data)
	return responseData, nil
}

func (osv *OrderService) CreateOrder(ctx context.Context, productId uuid.UUID, userId uuid.UUID) (model.Order, *shared.ErrorResponse) {
	productData, getErr := osv.productService.GetProductById(ctx, productId)
	if getErr != nil {
		return model.Order{}, getErr
	}

	if productData.Status != "active" {
		return model.Order{}, shared.NewErrorResponse(409, "you can't make order for this product. Product currently unavaible")
	}

	orderData, err := osv.repo.CreateOrders(ctx, userId, productId, productData.Price)
	if err != nil {
		return model.Order{}, shared.NewErrorResponse(500, "something wrong while creating new order try another time!")
	}

	return orderData, nil
}

func (osv *OrderService) GetAllOrderFromUser(ctx context.Context, userId uuid.UUID, status string) ([]OrderListDTO, *shared.ErrorResponse) {
	var tempStatus *string = nil
	if status != "" {
		tempStatus = &status
	}

	data, err := osv.repo.GetOrderFromUsers(ctx, userId, tempStatus)
	if err != nil {
		return []OrderListDTO{}, shared.NewErrorResponse(500, "something wrong while getting the order data! try again later")
	}

	responseData := orderListModelToDto(data)

	return responseData, nil
}

func (osv *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, status string) (model.Order, *shared.ErrorResponse) {
	data, err := osv.repo.UpdateOrderStatus(ctx, orderId, status)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, noOrderFound) {
			return model.Order{}, shared.NewErrorResponse(404, "no order with this id found!")
		}
		return model.Order{}, shared.NewErrorResponse(500, "something wrong while trying to update order status")
	}
	return data, nil
}

func (osv *OrderService) GetOrderById(ctx context.Context, id uuid.UUID) (model.Order, *shared.ErrorResponse) {
	data, err := osv.repo.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, noOrderFound) {
			return model.Order{}, shared.NewErrorResponse(404, "orders with this id not found!")
		}
	}

	return data, nil
}

func (osv *OrderService) DeleteOrders(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {

	err := osv.repo.DeleteOrder(ctx, id)
	if err != nil {

		if errors.Is(err, noOrderFound) {
			return shared.NewErrorResponse(404, "the order doesn't exist!")
		}

		return shared.NewErrorResponse(500, "something wrong while trying deleting the order")
	}

	return nil

}
