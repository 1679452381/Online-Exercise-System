package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type ProblemBasic struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`
	Content    string `gorm:"column:content;type:text;" json:"content"`
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11)" json:"max_runtime"` //最大运行时间
	MaxMem     int    `gorm:"column:max_mem;type:int(11)" json:"max_mem"`         //最大运行内存
}

func (ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword string) *gorm.DB {
	return utils.DB.Model(&ProblemBasic{}).Where("title LIKE ? OR content LIKE ? ", "%"+keyword+"%", "%"+keyword+"%")
}
