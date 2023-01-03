package routes

import "github.com/gorilla/mux"

// DEFINE ROUTES FOR API
func RouteInit(r *mux.Router) {
	AuthRoutes(r)
	UserRoutes(r)
	ProductRoutes(r)
	CartRoutes(r)
	TransactionRoutes(r)
}
