package service

import (
	"fmt"
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

	keyword := c.Query("keyword")
	offset := (page - 1) * size
	//	查数据库
	//count 记录数据的条数
	var count int64
	problems := make([]*models.Problem, 0)
	tx := models.GetProblemList(keyword)
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&problems).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"list": problems, "count": count}, c)
}
