package payment

import (
	"context"
	"errors"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	db *sqlcgen.Queries
}

func NewPaymentRepository(q *sqlcgen.Queries) *PaymentRepository {
	return &PaymentRepository{
		db: q,
	}
}

var ErrUserNotFound = errors.New("user not found");
var ErrOrderNotFound = errors.New("order not found");
var ErrPaymentNotFound = errors.New("payment not found");

func (Pr *PaymentRepository) CheckUserExist(ctx context.Context, userID uuid.UUID) (model.User,error){
	data,err := Pr.db.GetUserById(ctx,userID);
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return model.User{}, ErrUserNotFound;
		}
		return model.User{}, err;
	}
	return model.MapToUserModel(data),nil;
}

func (Pr *PaymentRepository) GetOrderDataById(ctx context.Context, userID, orderID uuid.UUID) (model.Order,error){
	data,err := Pr.db.GetOrderByID(ctx,orderID);
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return model.Order{}, ErrOrderNotFound;
		}
		return model.Order{}, err;
	}
	return model.OrderDataToModel(data.Order),nil;
}

func (Pr *PaymentRepository) GetLatestPayment(ctx context.Context, userID, orderID uuid.UUID)(model.Payment,error){
	
	data,err := Pr.db.GetLatestPayment(ctx,sqlcgen.GetLatestPaymentParams{
		UserID: userID,
		OrderID: orderID,
	})
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return model.Payment{}, ErrPaymentNotFound;
		}
		return model.Payment{}, err;
	}
	return model.MapToPaymentModel(data),nil;
}

func (Pr *PaymentRepository) CreatePayment(ctx context.Context, orderID uuid.UUID, method string, amount float64)(model.Payment,error){
	data,err := Pr.db.CreateNewPayment(ctx,sqlcgen.CreateNewPaymentParams{
		OrderID: orderID,
		Method: &method,
		Amount: amount,
	})
	if err != nil {
		return model.Payment{},err;
	}
	return model.MapToPaymentModel(data),nil;

}

func (Pr *PaymentRepository) UpdatePaymentWithGatewayData (ctx context.Context, oldPaymentData model.Payment, updateData PakasirResponse) (model.Payment,error){
	expiredAt,_ := time.Parse(time.RFC3339,updateData.Payment.ExpiredAt)
	data,err := Pr.db.UpdatePaymentWithGatewayData(ctx,sqlcgen.UpdatePaymentWithGatewayDataParams{
		ID: oldPaymentData.ID,
		Fee: &updateData.Payment.Fee,
		TotalPayment: &updateData.Payment.TotalPayment,
		ExpiredAt: &expiredAt,
	})
	if err != nil{
		return model.Payment{},err;
	}
	return model.MapToPaymentModel(data),nil;
}

func (Pr *PaymentRepository) UpdateExpired(ctx context.Context,paymentID uuid.UUID) error {
	err := Pr.db.SetPaymentExpired(ctx, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (Pr *PaymentRepository) GetPaymentByID(ctx context.Context, userID, paymentID uuid.UUID) (model.Payment,error){
	data, err := Pr.db.GetPaymentById(ctx,sqlcgen.GetPaymentByIdParams{
		ID: paymentID,
		UserID: userID,
	})
	if err != nil{
		if errors.Is(err,pgx.ErrNoRows){
			return model.Payment{},ErrPaymentNotFound;	
		}
			return model.Payment{},err;
	}
	return model.MapToPaymentModel(data),nil;
}

func (Pr *PaymentRepository) GetAllPaymentsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Payment,error){
	data, err := Pr.db.GetAllPayments(ctx,userID);
	if err != nil{
		return []model.Payment{},err;
	}
	return model.MapListToPaymentModel(data),nil;
}