package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

const (
	TypeOrderStatusPending    = "pending"
	TypeOrderStatusProcessing = "processing"
	TypeOrderStatusCompleted  = "completed"
	TypeOrderStatusCancelled  = "cancelled"
)

type Order struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	ProductID    uuid.UUID  `json:"product_id"`
	OrderedPrice float64    `json:"ordered_price"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	User    *User    `json:"user,omitempty"`
	Product *Product `json:"product,omitempty"`
}

func MapGenToOrder(gen sqlcgen.Order, u sqlcgen.User, p sqlcgen.Product) Order {
	user := MapToUserModel(u)
	product := MapToProductModel(p)
	return Order{
		ID:           gen.ID,
		UserID:       gen.UserID,
		ProductID:    gen.ProductID,
		OrderedPrice: gen.OrderedPrice,
		Status:       gen.Status,
		CreatedAt:    gen.CreatedAt,
		UpdatedAt:    gen.UpdatedAt,
		DeletedAt:    gen.DeletedAt,
		User:         &user,
		Product:      &product,
	}
}

func GenListToOrderMap(gen []sqlcgen.Order, u []sqlcgen.User, p []sqlcgen.Product) []Order {

	data := make([]Order, len(gen))
	for i := range gen {
		data[i] = MapGenToOrder(gen[i], u[i], p[i])
	}

	return data
}

func OrderDataToModel(gen sqlcgen.Order) Order {
	return Order{
		ID:           gen.ID,
		UserID:       gen.UserID,
		ProductID:    gen.ProductID,
		OrderedPrice: gen.OrderedPrice,
		Status:       gen.Status,
		CreatedAt:    gen.CreatedAt,
		UpdatedAt:    gen.UpdatedAt,
		DeletedAt:    gen.DeletedAt,
	}
}
