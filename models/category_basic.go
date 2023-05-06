package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:identity;type:varchar(100);" json:"name"` //分类名称
	ParentId int    `gorm:"column:category_id;type:int(11);" json:"parent_id" `
}

func (CategoryBasic) TableName() string {
	return "category_basic"
}

func GetCategoryList(keyword string) *gorm.DB {
	tx := utils.DB.Model(&CategoryBasic{})
	if keyword != "" {
		tx.Where("name like ?", "%"+keyword+"%")
	}
	return tx
}
