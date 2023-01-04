package database

import (
	"fmt"
	"server/models"
	"server/pkg/sql"
)

// CREATE DATABASE MYSQL
func RunMigration() {
	err := sql.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.Transaction{})
	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
