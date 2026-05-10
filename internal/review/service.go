package review

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ReviewService struct {
	repo *ReviewRepository
}

func NewReviewService(repo *ReviewRepository) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, userID uuid.UUID, req CreateReviewRequest) (model.Review, *shared.ErrorResponse) {
	// 1. Get order
	ord, err := s.repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Review{}, shared.NewErrorResponse(404, "order not found")
		}
		return model.Review{}, shared.NewErrorResponse(500, "failed to get order data")
	}

	// 2. Check ownership
	if ord.UserID != userID {
		return model.Review{}, shared.NewErrorResponse(403, "you are not authorized to review this order")
	}

	// 3. Check order status
	if ord.Status != "completed" {
		return model.Review{}, shared.NewErrorResponse(400, "only completed orders can be reviewed")
	}

	// 4. Check if review already exists
	_, err = s.repo.GetReviewByOrderID(ctx, req.OrderID)
	if err == nil {
		return model.Review{}, shared.NewErrorResponse(409, "this order has already been reviewed")
	}

	// 5. Get User and Product models
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return model.Review{}, shared.NewErrorResponse(500, "failed to get user data")
	}

	product, err := s.repo.GetProductByID(ctx, ord.ProductID)
	if err != nil {
		return model.Review{}, shared.NewErrorResponse(500, "failed to get product data")
	}

	// 6. Create review
	review, err := s.repo.CreateReview(ctx, user, product, ord, req.Rating, req.Comment)
	if err != nil {
		return model.Review{}, shared.NewErrorResponse(500, "failed to create review")
	}

	return review, nil
}

func (s *ReviewService) ListReviewsByProduct(ctx context.Context, productID uuid.UUID) ([]model.Review, *shared.ErrorResponse) {
	reviews, err := s.repo.ListReviewsByProduct(ctx, productID)
	if err != nil {
		return nil, shared.NewErrorResponse(500, "failed to get reviews")
	}
	return reviews, nil
}

func (s *ReviewService) ListFeaturedReviews(ctx context.Context, limit int32) ([]model.Review, *shared.ErrorResponse) {
	if limit <= 0 {
		limit = 5
	}
	reviews, err := s.repo.ListFeaturedReviews(ctx, limit)
	if err != nil {
		return nil, shared.NewErrorResponse(500, "failed to get featured reviews")
	}
	return reviews, nil
}

func (s *ReviewService) UpdateReviewShowStatus(ctx context.Context, id uuid.UUID, show bool) (model.Review, *shared.ErrorResponse) {
	review, err := s.repo.UpdateReviewShowStatus(ctx, id, show)
	if err != nil {
		return model.Review{}, shared.NewErrorResponse(500, "failed to update review status")
	}
	return review, nil
}
