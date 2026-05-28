package review

import (
	"strconv"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	svc *ReviewService
}

func NewReviewHandler(s *ReviewService) *ReviewHandler {
	return &ReviewHandler{
		svc: s,
	}
}

func (h *ReviewHandler) PostCreateReview(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(401, shared.NewErrorResponse(401, "unauthorized"))
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, err := h.svc.CreateReview(c, userID, req)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(201, shared.NewSuccessResponse(201, "Review successfully created", data))
}

func (h *ReviewHandler) HandleListReviewsByProduct(c *gin.Context) {
	productIDParams := c.Param("product_id")
	productID, err := uuid.Parse(productIDParams)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid product id"))
		return
	}

	data, getErr := h.svc.ListReviewsByProduct(c, productID)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "Reviews successfully retrieved", data))
}

func (h *ReviewHandler) HandleListFeaturedReviews(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		limit = 5
	}

	data, getErr := h.svc.ListFeaturedReviews(c, int32(limit))
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "Featured reviews successfully retrieved", data))
}

func (h *ReviewHandler) PatchUpdateReviewShowStatus(c *gin.Context) {
	idParams := c.Param("id")
	id, err := uuid.Parse(idParams)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid review id"))
		return
	}

	var req UpdateReviewShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, updErr := h.svc.UpdateReviewShowStatus(c, id, req.Show)
	if updErr != nil {
		c.JSON(updErr.Code, updErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "Review show status successfully updated", data))
}

func (h *ReviewHandler) HandleListMyReviews(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(401, shared.NewErrorResponse(401, "unauthorized"))
		return
	}

	data, err := h.svc.ListMyReviews(c, userID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "My reviews successfully retrieved", data))
}

func (h *ReviewHandler) RegisterRoutes(r gin.IRouter) {
	reviewApi := r.Group("/reviews")

	// Public routes
	reviewApi.GET("/product/:product_id", h.HandleListReviewsByProduct)
	reviewApi.GET("/featured", h.HandleListFeaturedReviews)

	// User routes
	userRoutes := reviewApi.Group("/")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.POST("/create", h.PostCreateReview)
	userRoutes.GET("/my-reviews", h.HandleListMyReviews)

	// Admin routes
	adminRoutes := reviewApi.Group("/")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
	adminRoutes.PATCH("/show-status/:id", h.PatchUpdateReviewShowStatus)
}
