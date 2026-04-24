package payment

type PaymentService struct {
	repo *PaymentRepository
}

func NewPaymentService(pr *PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: pr,
	}
}
