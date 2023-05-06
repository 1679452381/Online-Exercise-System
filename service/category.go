package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"online_exercise_system/utils"
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

// CategoryAdd
// @Tags 管理员私有方法
// @Summary 添加分类
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parent_identity formData string false "parent_identity"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/category_add [post]
func CategoryAdd(c *gin.Context) {
	name := c.PostForm("name")
	parentIdentity := c.PostForm("parent_identity")
	if name == "" {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	category := &models.CategoryBasic{
		Identity:       utils.GetUUID(),
		Name:           name,
		ParentIdentity: parentIdentity,
	}
	//查询数据库 看是否已有分类
	var count int64
	tx := models.CategoryBasicDB()
	err := tx.Where("name = ?", name).Count(&count).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	fmt.Println(count)
	if count != 0 {
		response.SuccessResponseWithMsg("已有该分类", c)
		return
	}
	//创建数据
	err = utils.DB.Create(&category).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	response.SuccessResponse("创建成功", category, c)
}

// CategoryDel
// @Tags 管理员私有方法
// @Summary 删除分类
// @Param authorization header string true "authorization"
// @Param category_identity query string true "category_identity"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/category_del [delete]
func CategoryDel(c *gin.Context) {
	categoryIdentity := c.Query("category_identity")
	if categoryIdentity == "" {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	var count int64
	err := utils.DB.Model(&models.ProblemCategory{}).Where("category_id=(select id from category_basic WHERE identity=? LIMIT 1)", categoryIdentity).Count(&count).Error
	if err != nil {
		response.FailResponseWithMsg("获取关联分类问题失败", c)
		return
	}
	if count > 0 {
		response.FailResponseWithMsg("该分类下已存在问题，不能删除", c)
		return
	}
	tx := models.CategoryBasicDB()
	err = tx.Where("identity=?", categoryIdentity).Or("parent_identity=?", categoryIdentity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	response.SuccessResponseWithMsg("删除成功", c)
}

// CategoryUpdate
// @Tags 管理员私有方法
// @Summary 更新分类
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param identity formData string true "identity"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/category_list [get]
func CategoryUpdate(c *gin.Context) {
	name := c.PostForm("name")
	identity := c.PostForm("identity")
	if name == "" || identity == "" {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	tx := models.CategoryBasicDB()
	category := &models.CategoryBasic{}
	err := tx.Model(&category).Where("identity=?", identity).Update("name", name).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	response.SuccessResponse("更新成功", category, c)
}
