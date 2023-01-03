package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/mysql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	// GET AUTH REPOSITORY HANDLER
	authRepository := repositories.RepositoryAuth(mysql.DB)
	h := handlers.HandlerAuth(authRepository)

	// DEFINE ROUTES
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/auth", middleware.Auth(h.CheckAuth)).Methods("GET")
}
