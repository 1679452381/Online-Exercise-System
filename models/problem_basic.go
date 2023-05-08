package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(100);" json:"identity"`
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"`
	Content           string             `gorm:"column:content;type:text;" json:"content"`
	MaxRuntime        int                `gorm:"column:max_runtime;type:int(11)" json:"max_runtime"` //最大运行时间
	MaxMem            int                `gorm:"column:max_mem;type:int(11)" json:"max_mem"`         //最大运行内存
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`                //关联分类表
	TestCases         []*TestCase        `gorm:"foreignKey:problem_identity;references:identity"`    //关联测试用例表
	PassNum           int                `gorm:"column:pass_num;type:int(11);"  json:"pass_num"`     // 通过问题次数
	SubmitNum         int                `gorm:"column:submit_num;type:int(11);"  json:"submit_num"` //提交次数
}

func (ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	// Preload 预加载
	tx := utils.DB.Model(&ProblemBasic{}).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("title LIKE ? OR content LIKE ? ", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id =(SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}

func GetProblemDetail(Identity string) *gorm.DB {
	// Preload 预加载
	return utils.DB.Model(&ProblemBasic{}).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("identity=? ", Identity)
}
