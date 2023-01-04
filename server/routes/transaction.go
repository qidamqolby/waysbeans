package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/sql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	// GET TRANSACTION REPOSITORY HANDLER
	transactionRepository := repositories.RepositoryTransaction(sql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	// DEFINE ROUTES
	r.HandleFunc("/admin/transaction", middleware.Auth(h.FindTransactions)).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.UpdateTransaction)).Methods("PATCH")
	r.HandleFunc("/user/transaction", middleware.Auth(h.GetUserTransactionByUserID)).Methods("GET")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
}
