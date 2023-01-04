package main

import (
	"fmt"
	"net/http"
	"os"
	"server/database"
	"server/pkg/sql"
	"server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// // GET ENVIRONMENT VARIABLES
	// errEnv := godotenv.Load()
	// if errEnv != nil {
	// 	panic("Failed to load env file")
	// }

	// DATABASE INITIALIZATION
	sql.DatabaseInit()

	// RUN MIGRATION
	database.RunMigration()

	// DEFINE ROUTE
	r := mux.NewRouter()

	// CREATE PREFIX SUBROUTER
	routes.RouteInit(r.PathPrefix("/waysbeans").Subrouter())

	// INITIALIZATION "uploads" FOLDER
	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))))

	// DEFINE ALLOWED CORS REQUEST
	var AllowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var AllowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"})
	var AllowedOrigins = handlers.AllowedOrigins([]string{"*"})

	var port = os.Getenv("PORT")

	fmt.Println("server running localhost:" + port)

	// RUN SERVER
	http.ListenAndServe(":"+port, handlers.CORS(AllowedHeaders, AllowedMethods, AllowedOrigins)(r))
}
