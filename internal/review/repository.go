package review

import (
	"context"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sqlcgen.Queries
}

func NewReviewRepository(db *sqlcgen.Queries) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, user model.User, product model.Product, order model.Order, rating int32, comment *string) (model.Review, error) {
	data, err := r.db.CreateReview(ctx, sqlcgen.CreateReviewParams{
		UserID:    user.ID,
		ProductID: product.ID,
		OrderID:   order.ID,
		Rating:    rating,
		Comment:   comment,
		Show:      true,
	})
	if err != nil {
		return model.Review{}, err
	}
	return model.MapToReviewModel(data), nil
}

func (r *ReviewRepository) GetReviewByOrderID(ctx context.Context, orderID uuid.UUID) (model.Review, error) {
	data, err := r.db.GetReviewByOrderID(ctx, orderID)
	if err != nil {
		return model.Review{}, err
	}
	return model.MapToReviewModel(data), nil
}

func (r *ReviewRepository) ListReviewsByProduct(ctx context.Context, productID uuid.UUID) ([]model.Review, error) {
	data, err := r.db.ListReviewsByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	return model.MapListToReviewWithUserModel(data), nil
}

func (r *ReviewRepository) ListFeaturedReviews(ctx context.Context, limit int32) ([]model.Review, error) {
	data, err := r.db.ListFeaturedReviews(ctx, limit)
	if err != nil {
		return nil, err
	}
	return model.MapListToFeaturedReviewModel(data), nil
}

func (r *ReviewRepository) UpdateReviewShowStatus(ctx context.Context, id uuid.UUID, show bool) (model.Review, error) {
	data, err := r.db.UpdateReviewShowStatus(ctx, sqlcgen.UpdateReviewShowStatusParams{
		ID:   id,
		Show: show,
	})
	if err != nil {
		return model.Review{}, err
	}
	return model.MapToReviewModel(data), nil
}

func (r *ReviewRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (model.Order, error) {
	data, err := r.db.GetOrderByID(ctx, id)
	if err != nil {
		return model.Order{}, err
	}
	return model.MapGenToOrder(data.Order, data.User, data.Product), nil
}

func (r *ReviewRepository) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	data, err := r.db.GetUserById(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return model.MapToUserModel(data), nil
}

func (r *ReviewRepository) GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	data, err := r.db.GetProductById(ctx, id)
	if err != nil {
		return model.Product{}, err
	}
	return model.MapToProductModel(data), nil
}

func (r *ReviewRepository) ListReviewsByUser(ctx context.Context, userID uuid.UUID) ([]model.Review, error) {
	data, err := r.db.ListReviewsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.MapListToReviewByUserModel(data), nil
}
