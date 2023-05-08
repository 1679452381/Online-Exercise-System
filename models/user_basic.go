package models

import (
	"gorm.io/gorm"
	"online_exercise_system/utils"
)

type UserBasic struct {
	gorm.Model
	Identity            string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	UserName            string `gorm:"column:username;type:varchar(100);" json:"user_name"`
	Phone               string `gorm:"column:phone;type:varchar(20);" json:"phone" `
	Password            string `gorm:"column:password;type:varchar(100);"json:"-"`
	Email               string `gorm:"column:email;type:varchar(100);"  json:"email"`
	CompletedProblemNum int    `gorm:"column:completed_problem_num;type:int(11);"  json:"completed_problem_num"`
	PassNum             int    `gorm:"column:pass_num;type:int(11);"  json:"pass_num"`     // 通过问题次数
	SubmitNum           int    `gorm:"column:submit_num;type:int(11);"  json:"submit_num"` //提交次数
	IsAdmin             int    `gorm:"column:is_admin;type:int(11);"  json:"is_admin"`
}

func (UserBasic) TableName() string {
	return "user_basic"
}

func GetUserBasicDetail(userIdentity string) *gorm.DB {
	return utils.DB.Model(&UserBasic{}).
		Where("identity=? ", userIdentity)
}

func GetUserBasicRankList() *gorm.DB {
	return utils.DB.Model(&UserBasic{})
}
