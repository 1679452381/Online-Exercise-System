package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity" `
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity;"` //关联问题基础表
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"` //关联用户基础表
	Path            string        `gorm:"column:path;type:varchar(255);" json:"path"`
	Status          int           `gorm:"column:status;type:int(1);" json:"status"` //0-待判断 1-答案正确 2-答案错误 3-运行超时 4-超出内存
}

func (SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity, status string) *gorm.DB {
	tx := utils.DB.Model(&SubmitBasic{}).Preload("ProblemBasic").Preload("UserBasic")
	if problemIdentity != "" {
		tx.Where("problem_identity=?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity=?", userIdentity)

	}
	if status != "" {
		tx.Where("status=?", status)

	}
	return tx
}
