package payment

import (
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
)

type PaymentRepository struct {
	db *sqlcgen.Queries
}

func NewPaymentRepository(q *sqlcgen.Queries) *PaymentRepository {
	return &PaymentRepository{
		db: q,
	}
}