package review

import "github.com/google/uuid"

type CreateReviewRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Rating  int32     `json:"rating" binding:"required,min=1,max=5"`
	Comment *string   `json:"comment"`
}

type UpdateReviewShowRequest struct {
	Show bool `json:"show" binding:"required"`
}
