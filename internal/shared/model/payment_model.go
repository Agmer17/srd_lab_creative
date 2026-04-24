package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)


type Payment struct{
	ID uuid.UUID `json:"payment_id"`
	OrderID uuid.UUID `json:"order_id"`
	Method *string `json:"method"`
	Status string `json:"status"`
	Amount float64 `json:"amount"`
	Fee float64 `json:"fee"`
	TotalPayment *float64 `json:"total_payment"`
	PaymentNumber *string `json:"payment_number"`
	ExpiredAt *time.Time `json:"expired_at"`
	PaidAt *time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
}

func MapToPaymentModel(data sqlcgen.Payment) Payment {
	var fee float64 = 0
	if data.Fee != nil {
		fee = *data.Fee
	}

	return Payment{
		ID:            data.ID,
		OrderID:       data.OrderID,
		Method:        data.Method,
		Status:        data.Status,
		Amount:        data.Amount,
		Fee:           fee,
		TotalPayment:  data.TotalPayment,
		PaymentNumber: data.PaymentNumber,
		ExpiredAt:     data.ExpiredAt,
		PaidAt:        data.PaidAt,
		CreatedAt:     data.CreatedAt,
	}
}

func MapListToPaymentModel(data []sqlcgen.Payment) []Payment {
	result := make([]Payment, len(data))
	for i, v := range data {
		result[i] = MapToPaymentModel(v)
	}
	return result
}