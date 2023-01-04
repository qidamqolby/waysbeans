package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/sql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router) {
	// GET PRODUCT REPOSITORY HANDLER
	productRepository := repositories.RepositoryProduct(sql.DB)
	h := handlers.HandlerProduct(productRepository)

	// DEFINE ROUTES
	r.HandleFunc("/products", h.FindProducts).Methods("GET")
	r.HandleFunc("/product/{id}", middleware.Auth(h.GetProduct)).Methods("GET")
	r.HandleFunc("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct))).Methods("POST")
	r.HandleFunc("/product/{id}", middleware.Auth(middleware.UploadFile(h.UpdateProduct))).Methods("PATCH")
	r.HandleFunc("/product/{id}", middleware.Auth(h.DeleteProduct)).Methods("DELETE")
}
