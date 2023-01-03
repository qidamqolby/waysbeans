package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/mysql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func CartRoutes(r *mux.Router) {
	// GET CART REPOSITORY HANDLER
	CartRepository := repositories.RepositoryCart(mysql.DB)
	h := handlers.HandlerCart(CartRepository)

	// DEFINE ROUTES
	r.HandleFunc("/cart", middleware.Auth(h.CreateCart)).Methods("POST")
	r.HandleFunc("/cart/{id}", middleware.Auth(h.DeleteCart)).Methods("DELETE")
	r.HandleFunc("/user/cart", middleware.Auth(h.FindCartByTransactionID)).Methods("GET")
}
