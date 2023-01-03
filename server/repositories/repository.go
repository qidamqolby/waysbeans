package repositories

import "gorm.io/gorm"

// DEFINE DATABASE FOR CRUD
type repository struct {
	db *gorm.DB
}
