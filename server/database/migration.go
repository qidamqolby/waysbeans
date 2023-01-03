package database

import (
	"fmt"
	"server/models"
	"server/pkg/mysql"
)

// CREATE DATABASE MYSQL
func RunMigration() {
	err := mysql.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.Transaction{})
	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
