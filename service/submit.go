package service

import (
	"github.com/gin-gonic/gin"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"strconv"
)

// SubmitList
// @Tags 公共方法
// @Summary 用户提交列表
// @Param page query string false "page"
// @Param size query string false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query  string false "user_identity"
// @Param status query  string false "status"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /submit/list [get]
func SubmitList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", global.DefaultPage))
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", global.DefaultSize))
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	problemIdentity := c.DefaultQuery("problem_identity", "")
	userIdentity := c.DefaultQuery("user_identity", "")
	status := c.DefaultQuery("status", "")

	offset := (page - 1) * size
	//	查数据库
	//count 记录数据的条数
	var count int64
	submits := make([]*models.SubmitBasic, 0)
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&submits).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	//fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"list": submits, "count": count}, c)
}
