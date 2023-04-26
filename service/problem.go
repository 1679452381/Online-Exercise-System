package service

import (
	"github.com/gin-gonic/gin"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"strconv"
)

// ProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query string false "page"
// @Param size query string false "size"
// @Param keyword query  string false "keyword"
// @Param category_identity query  string false "category_identity"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /problem/list [get]
func ProblemList(c *gin.Context) {
	//	获取page,size和keyword信息
	//  用DefaultQuery 在没有穿page 和size时 给默认值
	// strconv.Atoi()将字符串转化为int
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
	categoryIdentity := c.DefaultQuery("category_identity", "")

	keyword := c.Query("keyword")
	offset := (page - 1) * size
	//	查数据库
	//count 记录数据的条数
	var count int64
	problems := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&problems).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	//fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"list": problems, "count": count}, c)
}

// ProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param problem_identity query string false "problem_identity"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /problem/detail [get]
func ProblemDetail(c *gin.Context) {

	problemIdentity := c.DefaultQuery("problem_identity", "")

	//	查数据库
	problemDetail := &models.ProblemBasic{}
	tx := models.GetProblemDetail(problemIdentity)
	err := tx.Find(&problemDetail).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	//fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"list": problemDetail}, c)
}
