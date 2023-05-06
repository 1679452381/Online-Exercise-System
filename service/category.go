package service

import (
	"github.com/gin-gonic/gin"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"strconv"
)

// CategoryList
// @Tags 管理员私有方法
// @Summary 获取分类列表
// @Param authorization header string true "authorization"
// @Param page query string false "page"
// @Param size query string false "size"
// @Param keyword query  string false "keyword"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/category_list [get]
func CategoryList(c *gin.Context) {
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
	page = (page - 1) * size
	categories := make([]*models.CategoryBasic, 0)
	var count int64
	tx := models.GetCategoryList(keyword)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&categories).Error
	if err != nil {

		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	response.SuccessResponse("查询成功", gin.H{
		"data":  categories,
		"count": count,
	}, c)
}
