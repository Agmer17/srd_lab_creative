package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)


type Product struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description"`
	Price       float64    `json:"price"`
	Status      string     `json:"status"`
	IsFeatured  bool       `json:"is_featured"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func MapToProductModel(gr sqlcgen.Product) Product {
	return Product{
		ID:          gr.ID,
		Name:        gr.Name,
		Slug:        gr.Slug,
		Description: gr.Description,
		Price:		 gr.Price,
		Status:		 gr.Status,
		IsFeatured:  gr.IsFeatured,
		CreatedAt:   gr.CreatedAt,
		UpdatedAt:   gr.UpdatedAt,
	}
}

func MapListToProductModel(ls []sqlcgen.Product) []Product{

	tempList := make([]Product, len(ls))

	for i, v := range ls {
		tempList[i] = MapToProductModel(v)
	}

	return tempList
}
