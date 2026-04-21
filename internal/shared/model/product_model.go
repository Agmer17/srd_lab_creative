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

type ProductImage struct{
	ID uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	ImageUrl string `json:"image_url"`
	IsPrimary bool `json:"is_primary"`
	SortOrder int32 `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
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

// Memetakan 1 entitas gambar kembalian SQLC ke Model JSON
func MapToProductImageModel(img sqlcgen.ProductImage) ProductImage {
	return ProductImage{
		ID:        img.ID,
		ProductID: img.ProductID,
		ImageUrl:  img.ImageUrl,
		IsPrimary: img.IsPrimary,
		SortOrder: img.SortOrder, // Hapus fungsi pembungkus 'int()' ini jika di struct atas kamu ubah tipe SortOrder-nya jadi int32 sesuai saran!
		CreatedAt: img.CreatedAt,
	}
}


func MapToProductImageListModel(imgs []sqlcgen.ProductImage) []ProductImage {
	var list []ProductImage
	for _, img := range imgs {
		list = append(list, MapToProductImageModel(img))
	}
	return list
}

