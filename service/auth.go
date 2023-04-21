package service

import (
	"context"
	"fmt"
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
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","msg":"","token":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	//获取用户输入的账号密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		response.SuccessResponseWithMsg("账号或密码不能为空", c)
		return
	}
	//将密码进行md5加密并与数据库对比
	user := new(models.User)
	err := utils.DB.Where("username=? AND password = ?", username, utils.GetMd5(password)).First(&user).Error
	if err != nil {
		response.SuccessResponseWithMsg("用户名或密码错误", c)
		return
	}
	//生成token
	token, err := utils.GenerateToken(user.Identity, username)
	if err != nil {
		response.SuccessResponseWithMsg("系统错误", c)
		return
	}
	//将token存到redis中
	//fmt.Println("login", global.Token+username)
	err = utils.Redis.Set(context.Background(), global.Token+username, token, time.Hour*1).Err()
	if err != nil {
		response.SuccessResponseWithMsg("系统错误", c)
		return
	}
	//返回token
	response.SuccessResponseWithToken(token, c)
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Param email formData string false "email"
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /register [post]
func Register(c *gin.Context) {
	//	获取输入的信息
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	code := c.PostForm("code")
	if username == "" || password == "" || email == "" || code == "" {
		response.FailResponseWithMsg("参数有误", c)
		return
	}
	//  检测账号是否被注册
	user := &models.User{}
	count := utils.DB.Where("username = ?", username).Find(user).RowsAffected
	if count > 0 {
		response.SuccessResponseWithMsg("该账号已被使用", c)
		return
	}
	//	从redis中获取验证码 验证邮箱验证码
	//  创建用户
	u := models.User{
		Identity: utils.GetUUID(),
		UserName: username,
		Password: utils.GetMd5(password),
		Email:    email,
	}
	fmt.Println(u)
	err := utils.DB.Create(&u).Error
	if err != nil {
		response.FailResponseWithMsg("系统错误", c)
		return
	}
	fmt.Println(u)
	response.SuccessResponse("注册成功", u, c)
}
