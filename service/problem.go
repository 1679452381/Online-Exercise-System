package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"online_exercise_system/utils"
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

	problemIdentity := c.Query("problem_identity")
	if problemIdentity == "" {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
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

// CreateProblem
// @Tags 管理员私有方法
// @Summary 创建问题
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/problem_add [post]
func CreateProblem(c *gin.Context) {
	//获取参数
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(testCases) == 0 {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	//	创建问题
	problem := &models.ProblemBasic{
		Identity:   utils.GetUUID(),
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}
	//添加分类
	problemCategories := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		problemCategories = append(problemCategories, &models.ProblemCategory{
			ProblemId:  problem.ID,
			CategoryId: uint(categoryId),
		})
	}
	problem.ProblemCategories = problemCategories
	fmt.Println(problem.Identity)
	//	添加测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		//{"input":"1 2\n","input":"3\n"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			response.FailResponseWithMsg("测试用例格式错误", c)
			return
		}
		if _, ok := caseMap["input"]; !ok {
			response.FailResponseWithMsg("测试用例格式错误", c)
			return
		}
		if _, ok := caseMap["output"]; !ok {
			response.FailResponseWithMsg("测试用例格式错误", c)
			return
		}
		testCaseBasics = append(testCaseBasics, &models.TestCase{
			Identity:        utils.GetUUID(),
			ProblemIdentity: problem.Identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		})
	}
	problem.TestCases = testCaseBasics
	//数据库创建问题
	err := utils.DB.Create(problem).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回数据
	response.SuccessResponse("添加成功", gin.H{
		"problem_identity": problem.Identity,
	}, c)
}

// ModifyProblem
// @Tags 管理员私有方法
// @Summary 修改问题
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/admin/problem_modify [post]
func ModifyProblem(c *gin.Context) {
	//获取参数
	title := c.PostForm("title")
	problemIdentity := c.PostForm("identity")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if problemIdentity == "" || title == "" || content == "" || len(testCases) == 0 {
		response.FailResponseWithMsg("参数不能为空", c)
		return
	}
	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		//	问题基本信息保存
		problem := &models.ProblemBasic{
			Identity:   problemIdentity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMem:     maxMem,
		}
		err := tx.Where("identity=?", problemIdentity).Updates(problem).Error
		if err != nil {
			return err
		}
		//查询问题详情
		err = tx.Where("identity=?", problemIdentity).Find(&problem).Error
		if err != nil {
			return err
		}
		//关联问题分类更新
		//a.删除已存在的关联分类
		err = tx.Where("problem_id=?", problem.ID).Delete(&models.ProblemCategory{}).Error
		if err != nil {
			return err
		}
		//b.新增新的关联关系
		problemCategories := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			categoryId, _ := strconv.Atoi(id)
			problemCategories = append(problemCategories, &models.ProblemCategory{
				ProblemId:  problem.ID,
				CategoryId: uint(categoryId),
			})
		}
		err = tx.Create(&problemCategories).Error
		if err != nil {
			return err
		}

		//关联测试案例更新
		err = tx.Where("problem_identity=?", problem.Identity).Delete(&models.TestCase{}).Error
		if err != nil {
			return err
		}
		testCaseBasics := make([]*models.TestCase, 0)
		for _, testCase := range testCases {
			fmt.Println(testCase)
			//{"input":"1 2\n","input":"3\n"}
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {

				return errors.New("测试用例格式错误1")
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试用例格式错误2")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试用例格式错误3")
			}
			testCaseBasics = append(testCaseBasics, &models.TestCase{
				Identity:        utils.GetUUID(),
				ProblemIdentity: problem.Identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}

		err = tx.Create(testCaseBasics).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		response.FailResponseWithMsg(err.Error(), c)
		return
	}

	//	返回数据
	response.SuccessResponse("修改成功", gin.H{
		"problem_identity": problemIdentity,
	}, c)
}
