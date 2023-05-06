package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type CategoryBasic struct {
	gorm.Model
	Identity       string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	Name           string `gorm:"column:name;type:varchar(100);" json:"name"` //分类名称
	ParentIdentity string `gorm:"column:parent_identity;type:varchar(100);" json:"parent_identity" `
}

func (CategoryBasic) TableName() string {
	return "category_basic"
}

func CategoryBasicDB() *gorm.DB {
	return utils.DB.Model(&CategoryBasic{})
}

func GetCategoryList(keyword string) *gorm.DB {
	tx := utils.DB.Model(&CategoryBasic{})
	if keyword != "" {
		tx.Where("name like ?", "%"+keyword+"%")
	}
	return tx
}
