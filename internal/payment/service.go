package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (ps *PaymentService) GetGatewayPaymentStatus(ctx context.Context, paymentID uuid.UUID, amount float64) (PakasirStatusResponse, error) {
	reqUrl := fmt.Sprintf("https://app.pakasir.com/api/transactiondetail?project=%s&amount=%.0f&order_id=%s&api_key=%s",
		os.Getenv("depo_domain"), amount, paymentID.String(), os.Getenv("api_key"))

	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return PakasirStatusResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PakasirStatusResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PakasirStatusResponse{}, errors.New("Pakasir rejected our request with code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PakasirStatusResponse{}, err
	}

	var result PakasirStatusResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return PakasirStatusResponse{}, err
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
	if errors.Is(err, ErrPaymentNotFound) || latestPaymentData.Status == "cancelled"{
		fmt.Println("HIT");
		return ps.processNewPayment(ctx,orderData,paymentMethod);
	}

	// flow kalo ada.
	// flow udah dibayar
	if latestPaymentData.Status == "paid"{
		
		return model.Payment{}, shared.NewErrorResponse(400, "Payment has been paid");

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

func (ps *PaymentService) GetTransactionDetail (ctx context.Context, userID, paymentID uuid.UUID) (model.Payment, *shared.ErrorResponse){
	
	// verify id 
	_, err := ps.repo.CheckUserExist(ctx,userID);
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { 
			return model.Payment{}, shared.NewErrorResponse(404, "User does not exist")
		}
		
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the user data try again!")
	}
	// cari transaksinya
	data,err := ps.repo.GetPaymentByID(ctx,userID,paymentID);
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound){
			return model.Payment{}, shared.NewErrorResponse(404, "Payment does not exist");
		}
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while getting the payment data, try again!");
	}

	// lazy updatee
	if data.Status == "unpaid" && data.ExpiredAt != nil && time.Now().After(*data.ExpiredAt) {
		data.Status = "expired"
	}
	return data,nil;
}

func (ps *PaymentService) GetTransactionHistory (ctx context.Context, userID uuid.UUID) ([]model.Payment, *shared.ErrorResponse){
	// verify id real atau engga
	_, err := ps.repo.CheckUserExist(ctx,userID);
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { 
			return []model.Payment{}, shared.NewErrorResponse(404, "User does not exist")
		}
		
		return []model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the user data try again!")
	}
	// ambil data dari db
	rawData,getErr := ps.repo.GetAllPaymentsByUserID(ctx,userID);
	if getErr != nil {
		return []model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the payment data try again!");
	}
	
	// edit yang expired datanya
	for i := range rawData{
		if rawData[i].Status == "unpaid" && rawData[i].ExpiredAt != nil && time.Now().After(*rawData[i].ExpiredAt){
			rawData[i].Status = "expired";
		}
	}
	// return
	return rawData,nil;
}

func (ps *PaymentService) SyncTransaction (ctx context.Context, userID, paymentID uuid.UUID) (model.Payment, *shared.ErrorResponse){
	// verif user
	_, err := ps.repo.CheckUserExist(ctx,userID);
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { 
			return model.Payment{}, shared.NewErrorResponse(404, "User does not exist")
		}
		
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the user data try again!")
	}
	
	// cek payment idnya sesuai sama user id dan ada apa engga
	data,err := ps.repo.GetPaymentByID(ctx,userID,paymentID);
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound){
			return model.Payment{}, shared.NewErrorResponse(404, "Payment does not exist");
		}
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while getting the payment data, try again!");
	}
	
	// verif kalo payment yang diminta itu udah paid, expired, atau failed
	if data.Status == "paid" || data.Status == "expired" || data.Status == "failed" || (data.ExpiredAt != nil && time.Now().After(*data.ExpiredAt)){
		return model.Payment{}, shared.NewErrorResponse(400, "this payment cannot be synced again");
	}

	// call ke API
	responseData,err := ps.GetGatewayPaymentStatus(ctx,paymentID,data.Amount);
	if err != nil{
		return model.Payment{}, shared.NewErrorResponse(500, "Something went wrong while connecting to the payment gateway");
	}
	// liat resultnya terus compare masih sama apa berubah
	if responseData.Transaction.Status == "completed" && responseData.Transaction.CompletedAt != ""{
		// kalo berubah simpen ke DB terus return ke user
		responseData.Transaction.Status = "paid";
		updateData,err := ps.repo.UpdateWithResyncData(ctx,paymentID,responseData);
		if err != nil{
			return model.Payment{}, shared.NewErrorResponse(500, "Something went wrong while updating the payment data, try again");
		}
		return updateData,nil;
	}
	// kalo sama aja langsung return yang dari DB aja.
	return data,nil;
}

func (ps *PaymentService) CancelGatewayTransaction(ctx context.Context, paymentID uuid.UUID, amount float64) (PakasirStatusResponse, error) {
	// a. Siapkan payload (Bisa pakai PakasirRequest karena JSON bodynya sama persis)
	payload := PakasirRequest{
		Project: os.Getenv("depo_domain"),
		OrderID: paymentID,
		Amount:  amount,
		APIKey:  os.Getenv("api_key"),
	}

	jsonData, _ := json.Marshal(payload)
	reqUrl := "https://app.pakasir.com/api/transactioncancel"
	req, err := http.NewRequestWithContext(ctx, "POST", reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return PakasirStatusResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// b. Eksekusi request Cancel
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PakasirStatusResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PakasirStatusResponse{}, errors.New("Pakasir rejected cancel request with code: " + resp.Status)
	}

	// c. Baca & Unmarshal jawaban
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PakasirStatusResponse{}, err
	}

	// Karena format balasan cancel tidak memiliki root string "transaction",
	// kita pakai anonymous struct sementara untuk menampungnya agar tetap clean
	var result struct {
		Amount        float64 `json:"amount"`
		OrderID       string  `json:"order_id"`
		Project       string  `json:"project"`
		Status        string  `json:"status"`
		PaymentMethod string  `json:"payment_method"`
		CompletedAt   string  `json:"completed_at"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return PakasirStatusResponse{}, err
	}

	// d. Manipulasi hasil dan daur ulang struct luar kita (PakasirStatusResponse)
	finalRes := PakasirStatusResponse{}
	finalRes.Transaction.Amount = result.Amount
	finalRes.Transaction.OrderID = result.OrderID
	finalRes.Transaction.Project = result.Project
	finalRes.Transaction.Status = result.Status
	finalRes.Transaction.PaymentMethod = result.PaymentMethod
	finalRes.Transaction.CompletedAt = result.CompletedAt

	return finalRes, nil
}

func (ps *PaymentService) CancelPayment (ctx context.Context, userID, paymentID uuid.UUID) (model.Payment,*shared.ErrorResponse){
	// verif user
	_, err := ps.repo.CheckUserExist(ctx,userID);
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { 
			return model.Payment{}, shared.NewErrorResponse(404, "User does not exist")
		}
		
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while checking the user data try again!")
	}
	
	// cek payment idnya sesuai sama user id dan ada apa engga
	paymentData,err := ps.repo.GetPaymentByID(ctx,userID,paymentID);
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound){
			return model.Payment{}, shared.NewErrorResponse(404, "Payment does not exist");
		}
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while getting the payment data, try again!");
	}

	if paymentData.Status == "unpaid" && (paymentData.ExpiredAt != nil && time.Now().Before(*paymentData.ExpiredAt)){
		// call cancel
		_,cancelErr := ps.CancelGatewayTransaction(ctx,paymentData.ID,paymentData.Amount);
		if cancelErr != nil{
			return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while connecting to payment gateway");
		}

		// update ke DB
		updateData,updErr := ps.repo.UpdateCanceled(ctx,paymentData.ID);
		if updErr != nil{
			return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while updating the payment status");
		}

		// return hasil update-an DB.
		return updateData,nil;
	}
	// flow kalo yang diminta itu ternyata selain unpaid (expired, or anything)
	return model.Payment{}, shared.NewErrorResponse(400, "payments cannot be cancelled");
	
}

func (ps *PaymentService) WebHookVerification(ctx context.Context, req PakasirStatusResponse, paymentID uuid.UUID) (model.Payment, *shared.ErrorResponse){
	// verif paymentnya ada apa engga sekalian get data payment yang ada di DB.
	// cek payment idnya sesuai sama user id dan ada apa engga
	oldData,err := ps.repo.GetPaymentByIDWithoutUserID(ctx,paymentID);
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound){
			return model.Payment{}, shared.NewErrorResponse(404, "Payment does not exist");
		}
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while getting the payment data, try again!");
	}

	// call ke API pakasir
	responseData,err := ps.GetGatewayPaymentStatus(ctx,paymentID,oldData.Amount);
	if err != nil{
		return model.Payment{}, shared.NewErrorResponse(500, "Something went wrong while connecting to the payment gateway");
	}


	// bandingin
	// kalo beda
	if responseData.Transaction != req.Transaction {
		return model.Payment{}, shared.NewErrorResponse(403, "Invalid data, payment data is not authenticated");
	}

	// kalo sama update
	updData, errUpd := ps.repo.UpdateStatusFromWebhook(ctx, paymentID, req.Transaction.Status, req.Transaction.CompletedAt)
	if errUpd != nil {
		return model.Payment{}, shared.NewErrorResponse(500, "something went wrong while updating verified webhook data")
	}
	return updData, nil
}