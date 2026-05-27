package order

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type createOrderRequest struct {
	ProductId string `json:"product_id" binding:"required,uuid"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing completed cancelled"`
}

type OrderListDTO struct {
	ID           uuid.UUID             `json:"id"`
	UserID       uuid.UUID             `json:"user_id"`
	ProductID    uuid.UUID             `json:"product_id"`
	OrderedPrice float64               `json:"ordered_price"`
	Status       string                `json:"status"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	User         *orderUserDTO         `json:"user,omitempty"`
	Product      orderProductDTO       `json:"product"`
	Payment      []orderPaymentSummary `json:"payment"`
}

type orderUserDTO struct {
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
}

type orderPaymentSummary struct {
	ID     uuid.UUID `json:"payment_id"`
	Method *string   `json:"method"`
	Status string    `json:"status"`
	Amount float64   `json:"amount"`
}

type orderProductDTO struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func orderModelToDto(m model.Order) OrderListDTO {
	tempPaymentSummary := make([]orderPaymentSummary, len(m.Payment))
	for i, v := range m.Payment {
		tempPaymentSummary[i] = orderPaymentSummary{
			ID:     v.ID,
			Method: v.Method,
			Status: v.Status,
			Amount: v.Amount,
		}
	}

	return OrderListDTO{
		ID:           m.ID,
		Status:       m.Status,
		OrderedPrice: m.OrderedPrice,
		UserID:       m.UserID,
		ProductID:    m.ProductID,

		CreatedAt: m.CreatedAt,

		User: &orderUserDTO{
			FullName:       m.User.FullName,
			Email:          m.User.Email,
			ProfilePicture: m.User.ProfilePicture,
		},

		Product: orderProductDTO{
			Name: m.Product.Name,
			Slug: m.Product.Slug,
		},
		Payment: tempPaymentSummary,
	}
}

func orderListModelToDto(o []model.Order) []OrderListDTO {

	var data []OrderListDTO = make([]OrderListDTO, len(o))

	for i, v := range o {
		data[i] = orderModelToDto(v)
	}

	return data
}
