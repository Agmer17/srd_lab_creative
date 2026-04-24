package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

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

func (ps *PaymentService)RequestTransaction(ctx context.Context, paymentMethod string, paymentID uuid.UUID, amount float64)(PakasirResponse,error){
    // a. Siapkan Data JSON
	payload := PakasirRequest{
		Project: os.Getenv("depo_domain"),
		OrderID: paymentID,      
		Amount:  amount,
		APIKey:  os.Getenv("api_key"), 
	}
	
	jsonData, _ := json.Marshal(payload)
	reqUrl := "https://app.pakasir.com/api/transactioncreate/" + paymentMethod;
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return PakasirResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
    // c. Kirim Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PakasirResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PakasirResponse{}, errors.New("Pakasir rejected our requested with code: " + resp.Status)
	}
    // d. Baca Jawaban JSON
	bodyLengkap, err := io.ReadAll(resp.Body)
	if err != nil {
		return PakasirResponse{}, err
	}
    // e. Ubah Jawaban JSON menjadi struct Golang
	var result PakasirResponse
	err = json.Unmarshal(bodyLengkap, &result)
	if err != nil {
		return PakasirResponse{}, err
	}
	return result, nil	
}

func (ps *PaymentService) processNewPayment(ctx context.Context, orderData model.Order, paymentMethod string) (model.Payment, *shared.ErrorResponse) {
	// create DB
	newPaymentData, err := ps.repo.CreatePayment(ctx,orderData.ID,paymentMethod,orderData.OrderedPrice);
	if err != nil{
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while creating the payment data");
	}
	// call API
	gatewayData, err := ps.RequestTransaction(ctx,paymentMethod,newPaymentData.ID,newPaymentData.Amount);
	if err != nil{
		ps.repo.UpdateExpired(ctx,newPaymentData.ID);
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while connecting to the payment gateway");
	}
	// update DB
	finalData, err := ps.repo.UpdatePaymentWithGatewayData(ctx,newPaymentData,gatewayData);
	if err != nil{
		return model.Payment{}, shared.NewErrorResponse(500, "Something went wrong while updating the payment data");
	}

	// return datanya
	return finalData,nil; 
}


func (ps *PaymentService) AddTransaction (ctx context.Context, userID, orderID uuid.UUID, paymentMethod string) (model.Payment, *shared.ErrorResponse){

	// verify user id
	_, err := ps.repo.CheckUserExist(ctx,userID);
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { 
			return model.Payment{}, shared.NewErrorResponse(404, "User does not exist")
		}
		
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the user data try again!")
	}
	// verify order data
	orderData, err := ps.repo.GetOrderDataById(ctx, userID, orderID);
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) { 
			return model.Payment{}, shared.NewErrorResponse(404, "Order data does not exist");
		}
		
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the order data try again!");
	}

	// cek payment data ada atau tidak
	latestPaymentData, err := ps.repo.GetLatestPayment(ctx,userID,orderData.ID);

	// error handling untuk error diluar ga ada data payment
	if err != nil {
		if !(errors.Is(err, ErrPaymentNotFound)) { 
			return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the order data try again!");
		}
	}

	// ini kalo data payment ga ada (call ke pg kemudian simpen ke DB)
	if errors.Is(err, ErrPaymentNotFound){
		return ps.processNewPayment(ctx,orderData,paymentMethod);
	}

	// flow kalo ada.
	// flow udah dibayar
	if latestPaymentData.Status == "paid"{
		
		return model.Payment{}, shared.NewErrorResponse(400, "Pesanan ini sudah lunas");

	}else if latestPaymentData.ExpiredAt != nil && time.Now().After(*latestPaymentData.ExpiredAt) && latestPaymentData.Status == "unpaid"{
		// flow belum dibayar dan udah expired
		// update data biar expired
		err := ps.repo.UpdateExpired(ctx,latestPaymentData.ID);
		
		if err != nil{
			return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while expiring the payment data");
		}
		
		return ps.processNewPayment(ctx,orderData,paymentMethod);
	}else{
		// flow belumm dibayar dan belum expired
		return latestPaymentData, nil;
	}
	
}