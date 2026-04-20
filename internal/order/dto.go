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
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	OrderedPrice float64   `json:"ordered_price"`
	CreatedAt    time.Time `json:"created_at"`

	User    *orderUserDTO   `json:"user,omitempty"`
	Product orderProductDTO `json:"product"`
}

type orderUserDTO struct {
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
}

type orderProductDTO struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func orderModelToDto(m model.Order) OrderListDTO {
	return OrderListDTO{
		ID:           m.ID,
		Status:       m.Status,
		OrderedPrice: m.OrderedPrice,

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
	}
}

func orderListModelToDto(o []model.Order) []OrderListDTO {

	var data []OrderListDTO = make([]OrderListDTO, len(o))

	for i, v := range o {
		data[i] = orderModelToDto(v)
	}

	return data
}
