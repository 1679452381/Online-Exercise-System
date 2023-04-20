package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"online_exercise_system/global"
	"online_exercise_system/models"
	"online_exercise_system/response"
	"online_exercise_system/utils"
	"time"
)

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param account formData string false "account"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","msg":"","token":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	//获取用户输入的账号密码
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		response.SuccessResponseWithMsg("账号或密码不能为空", c)
		return
	}
	//将密码进行md5加密并与数据库对比
	user := new(models.User)
	err := utils.DB.Where("name=? AND password = ?", account, utils.GetMd5(password)).First(&user).Error
	if err != nil {
		response.SuccessResponseWithMsg("用户名或密码错误", c)
		return
	}
	//生成token
	token, err := utils.GenerateToken(user.Identity, account)
	if err != nil {
		response.SuccessResponseWithMsg("系统错误", c)
		return
	}
	//将token存到redis中
	utils.Redis.Set(context.Background(), global.Token+account, token, time.Hour*1)
	//返回token
	response.SuccessResponseWithToken(token, c)
}
