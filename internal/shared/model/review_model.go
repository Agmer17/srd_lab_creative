package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Rating    int32     `json:"rating"`
	Comment   *string   `json:"comment"`
	Show      bool      `json:"show"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    *User    `json:"user,omitempty"`
	Product *Product `json:"product,omitempty"`
}

func MapToReviewModel(gr sqlcgen.Review) Review {
	return Review{
		ID:        gr.ID,
		UserID:    gr.UserID,
		ProductID: gr.ProductID,
		OrderID:   gr.OrderID,
		Rating:    gr.Rating,
		Comment:   gr.Comment,
		Show:      gr.Show,
		CreatedAt: gr.CreatedAt,
		UpdatedAt: gr.UpdatedAt,
	}
}

func MapListToReviewModel(ls []sqlcgen.Review) []Review {
	tempList := make([]Review, len(ls))

	for i, v := range ls {
		tempList[i] = MapToReviewModel(v)
	}

	return tempList
}

func MapToReviewWithUserModel(r sqlcgen.Review, u sqlcgen.User) Review {
	review := MapToReviewModel(r)
	user := MapToUserModel(u)
	review.User = &user
	return review
}

func MapListToReviewWithUserModel(ls []sqlcgen.ListReviewsByProductRow) []Review {
	tempList := make([]Review, len(ls))

	for i, v := range ls {
		tempList[i] = MapToReviewWithUserModel(v.Review, v.User)
	}

	return tempList
}

func MapToFeaturedReviewModel(r sqlcgen.Review, u sqlcgen.User, p sqlcgen.Product) Review {
	review := MapToReviewWithUserModel(r, u)
	product := MapToProductModel(p)
	review.Product = &product
	return review
}

func MapListToFeaturedReviewModel(ls []sqlcgen.ListFeaturedReviewsRow) []Review {
	tempList := make([]Review, len(ls))

	for i, v := range ls {
		tempList[i] = MapToFeaturedReviewModel(v.Review, v.User, v.Product)
	}

	return tempList
}

func MapToReviewWithProductModel(r sqlcgen.Review, p sqlcgen.Product) Review {
	review := MapToReviewModel(r)
	product := MapToProductModel(p)
	review.Product = &product
	return review
}

func MapListToReviewByUserModel(ls []sqlcgen.ListReviewsByUserRow) []Review {
	tempList := make([]Review, len(ls))

	for i, v := range ls {
		tempList[i] = MapToReviewWithProductModel(v.Review, v.Product)
	}

	return tempList
}
