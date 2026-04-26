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

func (ph *PaymentHandler) PostCreatePayment(c *gin.Context) {
	// get user id
	userID, ok := middleware.GetUserID(c);
	if !ok {
		c.JSON(401,shared.NewErrorResponse(401,"Invalid session"));
		return;
	}
	// get order id
	path := c.Param("order_id");
	orderId, err := uuid.Parse(path);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return;
	}
	// get query params untuk method pembayaran
	paymentMethod := c.Query("method");
	
	// Validasi apakah query kosong
	if paymentMethod == "" {
		c.JSON(400, shared.NewErrorResponse(400, "Please input method query parameter. Example: ?method=qris"));
		return
	}
	// Validasi Enum (Apakah value yang diketik terdaftar di map ValidPaymentMethods?)
	if !ValidPaymentMethods[paymentMethod] {
		c.JSON(400, shared.NewErrorResponse(400, "Payment method invalid or not supported"));
		return
	}
	// langsung execute logic createTransaction
	data,createErr := ph.svc.AddTransaction(c,userID,orderId,paymentMethod);
	if createErr != nil{
		c.JSON(createErr.Code,createErr);
		return;
	}
	c.JSON(200, shared.NewSuccessResponse(200,"Payment successfully created",data));
	return;

}

func (ph *PaymentHandler) PostWebhookListener(c *gin.Context) {
	// Endpoint yang dipanggil oleh Payment Gateway ketika ada konfirmasi mutasi bayar (sukses/gagal/kadaluarsa).
}

func (ph *PaymentHandler) GetPaymentDetail(c *gin.Context) {
	// Menampilkan rincian dari DB lokal.
	// get user id
	userID, ok := middleware.GetUserID(c);
	if !ok {
		c.JSON(401,shared.NewErrorResponse(401,"Invalid session"));
		return;
	}
	// get payment id
	path := c.Param("payment_id");
	paymentID, err := uuid.Parse(path);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return;
	}

	// get data transaksi
	data, getErr := ph.svc.GetTransactionDetail(c,userID,paymentID);
	if getErr != nil{
		c.JSON(getErr.Code,getErr);
		return;
	}

	
	c.JSON(200, shared.NewSuccessResponse(200,"Payment successfully retrieved",data));
	return;

}

func (ph *PaymentHandler) PostCancelPayment(c *gin.Context) {
	// Membatalkan transaksi di Payment Gateway dan melakukan soft-delete (deleted_at + status canceled) di DB lokal.
}

func (ph *PaymentHandler) HandleGetPaymentHistory(c *gin.Context) {
	// Menampilkan list history transaksi user.
	
	// get user id
	userID, ok := middleware.GetUserID(c);
	if !ok {
		c.JSON(401,shared.NewErrorResponse(401,"Invalid session"));
		return;
	}

	// get data
	data, errGet := ph.svc.GetTransactionHistory(c,userID);
	if errGet != nil{
		c.JSON(errGet.Code,errGet);
		return;
	}
	
	c.JSON(200,shared.NewSuccessResponse(200,"Payment history successfully retrieved", data));
}

func (ph *PaymentHandler) PostManualSync(c *gin.Context) {
	// Mengambil status real-time dari API Payment Gateway lalu melakukan pembaruan di DB lokal apabila webhook meleset.
	
	// get user id
	userID, ok := middleware.GetUserID(c);
	if !ok {
		c.JSON(401,shared.NewErrorResponse(401,"Invalid session"));
		return;
	}
	
	// get payment id
	path := c.Param("payment_id");
	paymentID, err := uuid.Parse(path);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return;
	}

	// manual sync
	data,syncErr := ph.svc.SyncTransaction(c,userID,paymentID);
	if syncErr != nil{
		c.JSON(syncErr.Code,syncErr);
		return;
	}

	c.JSON(200, shared.NewSuccessResponse(200,"Payments data retrieved",data));
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
	protectedPaymentApi.POST("/create/:order_id", ph.PostCreatePayment);

	// user idnya diambil dari middleware aja terus dicek
	// Transaction History
	protectedPaymentApi.GET("/history", ph.HandleGetPaymentHistory);

	// Get Transaction Detail
	protectedPaymentApi.GET("/detail/:payment_id", ph.GetPaymentDetail);

	// Cancel Transaction
	protectedPaymentApi.POST("/cancel/:payment_id", ph.PostCancelPayment);

	// Manual Sync
	protectedPaymentApi.POST("/sync/:payment_id", ph.PostManualSync);

	// Webhook Listener
	paymentApi.POST("/webhook", ph.PostWebhookListener);
}