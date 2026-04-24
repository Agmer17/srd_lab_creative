package payment

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	svc *PaymentService
}

func NewPaymentHandler(s *PaymentService) *PaymentHandler {
	return &PaymentHandler{
		svc: s,
	}
}

func (ph *PaymentHandler) PostCreateTransaction(c *gin.Context) {
	// get user id
	userID, ok := middleware.GetUserID(c);
	if !ok {
		c.JSON(401,shared.NewErrorResponse(401,"Invvalid session"));
		return;
	}
	// get order id
	path := c.Param("order_id");
	orderId, err := uuid.Parse(path);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return;
	}	
	// langsung execute logic createTransaction
	data,createErr := ph.svc.AddTransaction(c,userID,orderId);
	if createErr != nil{
		c.JSON(createErr.Code,createErr);
		return;
	}
	c.JSON(200, shared.NewSuccessResponse(200,"Transaction successfully created",data));

}

func (ph *PaymentHandler) PostWebhookListener(c *gin.Context) {
	// Endpoint yang dipanggil oleh Payment Gateway ketika ada konfirmasi mutasi bayar (sukses/gagal/kadaluarsa).
}

func (ph *PaymentHandler) GetTransactionDetail(c *gin.Context) {
	// Menampilkan rincian dari DB lokal.
}

func (ph *PaymentHandler) PostCancelTransaction(c *gin.Context) {
	// Membatalkan transaksi di Payment Gateway dan melakukan soft-delete (deleted_at + status canceled) di DB lokal.
}

func (ph *PaymentHandler) HandleGetTransactionHistory(c *gin.Context) {
	// Menampilkan list history transaksi user.
}

func (ph *PaymentHandler) PostManualSync(c *gin.Context) {
	// Mengambil status real-time dari API Payment Gateway lalu melakukan pembaruan di DB lokal apabila webhook meleset.
}

func (ph *PaymentHandler) RegisterRoutes(r gin.IRouter) {
	paymentApi := r.Group("/payment");
	protectedPaymentApi := paymentApi.Group("/");
	protectedPaymentApi.Use(middleware.AuthMiddleware());
	protectedPaymentApi.Use(middleware.RoleMiddleware(middleware.RoleUser));


	/*
	Buat ngerequest ke transaction PG
	// order_id di id params
	{
    "project": "depodomain",
    "order_id": "INV123123",
    "amount": 99000,
    "api_key": "xxx123"
	}
	*/

	// Create Transaction
	protectedPaymentApi.POST("/create/:order_id", ph.PostCreateTransaction);

	// user idnya diambil dari middleware aja terus dicek
	// Transaction History
	protectedPaymentApi.GET("/history", ph.HandleGetTransactionHistory);

	// Get Transaction Detail
	protectedPaymentApi.GET("/detail/:payment_id", ph.GetTransactionDetail);

	// Cancel Transaction
	protectedPaymentApi.POST("/cancel/:payment_id", ph.PostCancelTransaction);

	// Manual Sync
	protectedPaymentApi.POST("/sync/:payment_id", ph.PostManualSync);

	// Webhook Listener
	paymentApi.POST("/webhook", ph.PostWebhookListener);
}