package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:identity;type:varchar(100);" json:"name"`
	Phone    string `gorm:"column:category_id;type:varchar(20);" json:"phone" `
	Password string `gorm:"column:title;type:varchar(32);" json:"password"`
	Email    string `gorm:"column:content;type:varchar(100);"  json:"email"`
}

func (UserModel) UserModelTableName() string {
	return "user"
}
