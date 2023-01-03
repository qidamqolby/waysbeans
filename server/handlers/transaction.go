package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	dto "server/dto/result"
	transactiondto "server/dto/transaction"
	"server/models"
	"server/repositories"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ID FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET REQUEST AND DECODING JSON
	request := new(transactiondto.TransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transaction, err := h.TransactionRepository.GetTransactionByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Cart Failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK UPDATE VALUE
	if request.Name != "" {
		transaction.Name = request.Name
	}

	if request.Email != "" {
		transaction.Email = request.Email
	}

	if request.Phone != "" {
		transaction.Phone = request.Phone
	}

	if request.Address != "" {
		transaction.Address = request.Address
	}

	transaction.Status = "pending"
	transaction.Total = request.Total
	transaction.UpdateAt = time.Now()

	// RUN REPOSITORY UPDATE TRANSACTION
	_, err = h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// SETUP FOR MIDTRANS
	DataSnap, _ := h.TransactionRepository.GetTransactionNotification(int(transaction.ID))

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(DataSnap.ID)),
			GrossAmt: int64(DataSnap.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: DataSnap.User.Name,
			Email: DataSnap.User.Email,
		},
	}

	// RUN MIDTRANS SNAP
	snapResp, _ := s.CreateTransaction(req)

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "Success", Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	// CHECK ROLE ADMIN
	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "You're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY FIND TRANSACTIONS
	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetUserTransactionByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ID FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transactions, err := h.TransactionRepository.GetUserTransactionByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: transactions}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderID := notificationPayload["order_id"].(string)

	transaction, _ := h.TransactionRepository.GetTransactionMidtrans(orderID)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			h.TransactionRepository.UpdateTransactionMidtrans("pending", int(transaction.ID))
		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			SendMail("success", transaction)
			h.TransactionRepository.UpdateTransactionMidtrans("success", int(transaction.ID))
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransactionMidtrans("success", int(transaction.ID))
	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		SendMail("failed", transaction)
		h.TransactionRepository.UpdateTransactionMidtrans("failed", int(transaction.ID))
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		SendMail("failed", transaction)
		h.TransactionRepository.UpdateTransactionMidtrans("waiting", int(transaction.ID))
	} else if transactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		h.TransactionRepository.UpdateTransactionMidtrans("waiting", int(transaction.ID))
	}

	w.WriteHeader(http.StatusOK)
}

func SendMail(status string, transaction models.Transaction) {

	if status != transaction.Status && (status == "success") {
		// GET VARIABLES FROM ENV
		var CONFIG_SMTP_HOST = os.Getenv("HOST_SYSTEM")
		var CONFIG_SMTP_PORT, _ = strconv.Atoi(os.Getenv("PORT_SYSTEM"))
		var CONFIG_SENDER_NAME = os.Getenv("SENDER_SYSTEM")
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

		var productName = transaction.Cart
		var price = strconv.Itoa(int(transaction.Total))

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", transaction.User.Email)
		mailer.SetHeader("Subject", "Transaction Status")
		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	  <html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<style>
		  h1 {
		  color: brown;
		  }
		</style>
		</head>
		<body>
		<h2>Product payment :</h2>
		<ul style="list-style-type:none;">
		  <li>Name : %v</li>
		  <li>Total payment: Rp.%v</li>
		  <li>Status : <b>%v</b></li>
		</ul>
		</body>
	  </html>`, productName, price, status))

		dialer := gomail.NewDialer(
			CONFIG_SMTP_HOST,
			CONFIG_SMTP_PORT,
			CONFIG_AUTH_EMAIL,
			CONFIG_AUTH_PASSWORD,
		)

		err := dialer.DialAndSend(mailer)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Mail sent! to " + transaction.User.Email)
	}
}
