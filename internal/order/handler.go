package order

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	svc *OrderService
}

func NewOrderHandler(sv *OrderService) *OrderHandler {
	return &OrderHandler{
		svc: sv,
	}
}

func (oh *OrderHandler) HandleGetAllOrders(c *gin.Context) {
	query := c.Query("status")
	data, err := oh.svc.GetAllOrders(c.Request.Context(), query)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "sucessfully getting the order data", data))

}

func (oh *OrderHandler) HandleGetOrderFromUsers(c *gin.Context) {

	id, _ := middleware.GetUserID(c)
	query := c.Query("status")

	data, err := oh.svc.GetAllOrderFromUser(c.Request.Context(), id, query)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "sucessfully getting the order data", data))
}

func (oh *OrderHandler) PostCreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	userID, _ := middleware.GetUserID(c)

	productUUID, errParse := uuid.Parse(req.ProductId)
	if errParse != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid product id"))
		return
	}

	data, err := oh.svc.CreateOrder(c.Request.Context(), productUUID, userID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(201, shared.NewSuccessResponse(201, "order created successfully", data))
}

func (oh *OrderHandler) HandleGetOrderByID(c *gin.Context) {
	idParam := c.Param("id")

	orderID, errParse := uuid.Parse(idParam)
	if errParse != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid order id"))
		return
	}

	data, err := oh.svc.GetOrderById(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting order detail", data))
}

func (oh *OrderHandler) PatchUpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")

	orderID, errParse := uuid.Parse(idParam)
	if errParse != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid order id"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, err := oh.svc.UpdateOrderStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "order status updated successfully", data))
}

func (oh *OrderHandler) HandleDeleteOrder(c *gin.Context) {
	idParam := c.Param("id")

	orderID, errParse := uuid.Parse(idParam)
	if errParse != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid order id"))
		return
	}

	err := oh.svc.DeleteOrders(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "order deleted successfully", nil))
}

func (oh *OrderHandler) RegisterRoutes(r gin.IRouter) {

	orderApi := r.Group("/orders")
	orderApi.Use(middleware.AuthMiddleware())

	// user routes
	orderApi.GET("/my-orders", oh.HandleGetOrderFromUsers)
	orderApi.POST("/create", oh.PostCreateOrder)
	orderApi.GET("/:id", oh.HandleGetOrderByID)

	// admin routes
	adminOrder := orderApi.Group("/")
	adminOrder.Use(middleware.RoleMiddleware(middleware.RoleAdmin))

	adminOrder.GET("/get-all", oh.HandleGetAllOrders)
	adminOrder.PATCH("/:id/status", oh.PatchUpdateOrderStatus)
	adminOrder.DELETE("/:id", oh.HandleDeleteOrder)
}
