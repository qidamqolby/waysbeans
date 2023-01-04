package routes

import (
	"server/handlers"
	"server/pkg/middleware"
	"server/pkg/sql"
	"server/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	// GET USER REPOSITORY HANDLER
	userRepository := repositories.RepositoryUser(sql.DB)
	h := handlers.HandlerUser(userRepository)

	//DEFINE ROUTES
	r.HandleFunc("/user", middleware.Auth(h.GetUser)).Methods("GET")
	r.HandleFunc("/user", middleware.Auth(middleware.UploadFile(h.UpdateUser))).Methods("PATCH")

}
