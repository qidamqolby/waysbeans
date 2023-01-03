package models

import "time"

type Transaction struct {
	ID       int32        `json:"id" gorm:"primary_key: auto_increment"`
	UserID   int32        `json:"user_id"`
	User     UserResponse `json:"user"`
	Name     string       `json:"name" form:"name" gorm:"type: varchar(50)"`
	Email    string       `json:"email" gorm:"type: varchar(50)"`
	Phone    string       `json:"phone" form:"phone" gorm:"type: varchar(50)"`
	Address  string       `json:"address" form:"address" gorm:"type : text"`
	Status   string       `json:"status" gorm:"type: varchar(50)"`
	Total    int32        `json:"total" gorm:"type: int"`
	Cart     []Cart       `json:"cart"`
	CreateAt time.Time    `json:"-"`
	UpdateAt time.Time    `json:"update_at"`
}
