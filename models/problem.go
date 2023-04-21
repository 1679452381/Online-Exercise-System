package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type Problem struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	CategoryId string `gorm:"column:category_id;type:varchar(255);" json:"category_id "`
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`
	Content    string `gorm:"column:content;type:text;" json:"content"`
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11)" json:"max_runtime"` //最大运行时间
	MaxMem     int    `gorm:"column:max_mem;type:int(11)" json:"max_mem"`         //最大运行内存
}

func (Problem) TableName() string {
	return "problem"
}

func GetProblemList(keyword string) *gorm.DB {
	return utils.DB.Model(&Problem{}).Where("title LIKE ? OR content LIKE ? ", "%"+keyword+"%", "%"+keyword+"%")
}
