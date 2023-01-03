package handlers

import (
	"encoding/json"
	"net/http"
	cartdto "server/dto/cart"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// SETUP HANLDER STRUCT
type handlerCart struct {
	CartRepository repositories.CartRepository
}

// SETUP HANDLER FUNCTION
func HandlerCart(CartRepository repositories.CartRepository) *handlerCart {
	return &handlerCart{CartRepository}
}

// CREATE CART
func (h *handlerCart) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET REQUEST AND DECODING JSON
	cartRequest := new(cartdto.CartRequest)
	if err := json.NewDecoder(r.Body).Decode(&cartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY GET PRODUCT BY PRODUCT ID
	product, err := h.CartRepository.GetProductCartByID(int(cartRequest.ProductID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// FIND TOTAL PRICE PRODUCT FROM QUANTITY REQUEST
	total := product.Price * cartRequest.Quantity

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	userTransaction, err := h.CartRepository.GetCartTransactionByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK IF EXIST
	if userTransaction.ID == 0 {

		// SETUP FOR QUERY TRANSACTION
		transaction := models.Transaction{
			ID:       int32(time.Now().Unix()),
			UserID:   int32(userID),
			Status:   "waiting",
			Total:    0,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}

		// RUN REPOSITORY CREATE TRANSACTION
		transactionData, err := h.CartRepository.CreateTransaction(transaction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		// SETUP FOR QUERY CART
		cart := models.Cart{
			UserID:        int32(userID),
			ProductID:     cartRequest.ProductID,
			Product:       models.Product{},
			OrderQty:      cartRequest.Quantity,
			Subtotal:      total,
			TransactionID: transactionData.ID,
			CreateAt:      time.Now(),
		}

		// RUN REPOSITORY CREATE CART
		data, err := h.CartRepository.CreateCart(cart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataResponse, _ := h.CartRepository.GetCart(int(data.ID))

		// WRITE RESPONSE
		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: "success", Data: dataResponse}
		json.NewEncoder(w).Encode(response)
	} else {

		// SETUP FOR QUERY CART
		cart := models.Cart{
			UserID:        int32(userID),
			ProductID:     cartRequest.ProductID,
			Product:       models.Product{},
			OrderQty:      cartRequest.Quantity,
			Subtotal:      total,
			TransactionID: userTransaction.ID,
			CreateAt:      time.Now(),
		}

		// RUN REPOSITORY CREATE CART
		data, err := h.CartRepository.CreateCart(cart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataResponse, _ := h.CartRepository.GetCart(int(data.ID))

		// WRITE RESPONSE
		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: "success", Data: dataResponse}
		json.NewEncoder(w).Encode(response)
	}
}

// DELETE CART
func (h *handlerCart) DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET CART ID FROM URL
	CartID, _ := strconv.Atoi(mux.Vars(r)["id"])

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET CART
	cart, err := h.CartRepository.GetCart(CartID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// VALIDATE REQUEST BY USER ID
	if userID != int(cart.UserID) {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// DELETE DATA
	data, err := h.CartRepository.DeleteCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION FIND CARTS BY TRANSACTION ID
func (h *handlerCart) FindCartByTransactionID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transaction, err := h.CartRepository.GetCartTransactionByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY FIND CARTS BY TRANSACTION ID
	carts, err := h.CartRepository.FindCartByTransactionID(int(transaction.ID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: carts}
	json.NewEncoder(w).Encode(response)
}
