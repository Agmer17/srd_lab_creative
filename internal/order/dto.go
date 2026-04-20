package order

import (
	"time"

	"github.com/google/uuid"
)

type createOrderRequest struct {
	ProductId string `json:"product_id" binding:"required,uuid"`
}

type updateOrderStatusRequest struct {
	Id     string `json:"order_id"`
	Status string `json:"order_status"`
}

type orderListDTO struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	OrderedPrice float64   `json:"ordered_price"`
	CreatedAt    time.Time `json:"created_at"`

	User    *OrderUserDTO   `json:"user,omitempty"`
	Product OrderProductDTO `json:"product"`
}

type OrderUserDTO struct {
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
}

type OrderProductDTO struct {
	Name string `json:"name"`
}
