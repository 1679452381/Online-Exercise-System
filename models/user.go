package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	UserName string `gorm:"column:username;type:varchar(100);" json:"user_name"`
	Phone    string `gorm:"column:phone;type:varchar(20);" json:"phone" `
	Password string `gorm:"column:password;type:varchar(100);" json:"_"`
	Email    string `gorm:"column:email;type:varchar(100);"  json:"email"`
}

func (User) TableName() string {
	return "users"
}
