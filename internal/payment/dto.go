package payment

import "github.com/google/uuid"

type PaymentMethod string

const (
	MethodCimbNiagaVA  PaymentMethod = "cimb_niaga_va"
	MethodBNIVA        PaymentMethod = "bni_va"
	MethodQris         PaymentMethod = "qris"
	MethodSampoernaVA  PaymentMethod = "sampoerna_va"
	MethodBNCVA        PaymentMethod = "bnc_va"
	MethodMaybankVA    PaymentMethod = "maybank_va"
	MethodPermataVA    PaymentMethod = "permata_va"
	MethodAtmBersamaVA PaymentMethod = "atm_bersama_va"
	MethodArthaGrahaVA PaymentMethod = "artha_graha_va"
	MethodBRIVA        PaymentMethod = "bri_va"
)

var ValidPaymentMethods = map[string]bool{
	string(MethodCimbNiagaVA):  true,
	string(MethodBNIVA):        true,
	string(MethodQris):         true,
	string(MethodSampoernaVA):  true,
	string(MethodBNCVA):        true,
	string(MethodMaybankVA):    true,
	string(MethodPermataVA):    true,
	string(MethodAtmBersamaVA): true,
	string(MethodArthaGrahaVA): true,
	string(MethodBRIVA):        true,
}
type PakasirRequest struct {
	Project string  `json:"project"`
	OrderID uuid.UUID  `json:"order_id"`
	Amount  float64 `json:"amount"`
	APIKey  string  `json:"api_key"`
}

type PakasirResponse struct {
	Payment struct {
		Project       string  `json:"project"`
		OrderID       string  `json:"order_id"`
		Amount        float64 `json:"amount"`
		Fee           float64 `json:"fee"`
		TotalPayment  float64 `json:"total_payment"`
		PaymentMethod string  `json:"payment_method"`
		PaymentNumber string  `json:"payment_number"`
		ExpiredAt     string  `json:"expired_at"` // Asumsinya pakasi me-return string RFC3339
	} `json:"payment"`
}

type PakasirStatusResponse struct {
	Transaction struct {
		Amount        float64 `json:"amount" binding:"required"`
		OrderID       string  `json:"order_id" binding:"required"`
		Project       string  `json:"project"`
		Status        string  `json:"status" binding:"required"`
		PaymentMethod string  `json:"payment_method"`
		CompletedAt   string  `json:"completed_at"`
	} `json:"transaction"`
}