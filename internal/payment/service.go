package payment

import (
	"context"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type PaymentService struct {
	repo *PaymentRepository
}

func NewPaymentService(pr *PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: pr,
	}
}


func (ps *PaymentService) AddTransaction (ctx context.Context, userID, orderID uuid.UUID) (model.Payment, shared.ErrorResponse){
	// Happy path
	// Verify User ID
	// Verify OrderID
	// Cek apakah ada transaction ID yang connect ke User ID dan OrderID tersebut
	// Ga ada
	// Generate Call API
	// Save ke DB
	// Return ke User

	// Udah kegenerate Path
	// Verify User ID
	// Verify OrderID
	// Cek apakah ada transaction ID yang connect ke User ID dan OrderID tersebut
	// Cek Expiry Date
	// not expired
	// return data yang di transaction ID

	// Udah kegenerate tapi expired
	// Verify User ID
	// Verify Order ID
	// Cek apakah ada transaction ID yang connect ke User ID dan OrderID tersebut
	// Cek Expiry Date
	// expired
	// Update datanya jadi expired sekaligus soft delete
	// Call ke API lagi (orderID kita sedikit modified tambahin orderID-2 atau -3 Pas kita call)
	// Save ke DB yang baru
	// Return ke User

	
}