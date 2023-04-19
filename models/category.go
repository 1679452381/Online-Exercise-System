package models

import "gorm.io/gorm"

type CategoryModel struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:identity;type:varchar(100);" json:"name"` //分类名称
	ParentId int    `gorm:"column:category_id;type:int(11);" json:"parent_id" `
}

func (CategoryModel) CategoryModelTableName() string {
	return "category"
}
