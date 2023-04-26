package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	user := new(models.UserBasic)
	err := utils.DB.Where("username=? AND password = ?", username, utils.GetMd5(password)).First(&user).Error
	if err != nil {
		log.Fatalln(err)
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
// @Param code formData string false "code"
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
	user := &models.UserBasic{}
	count := utils.DB.Where("username = ?", username).Find(user).RowsAffected
	if count > 0 {
		response.SuccessResponseWithMsg("该账号已被使用", c)
		return
	}
	fmt.Println(code)
	//	从redis中获取验证码 验证邮箱验证码
	emailCode, err := utils.Redis.Get(context.Background(), global.EmailCode+email).Result()
	fmt.Println(code, emailCode)
	if err != nil {
		log.Printf("%v", err.Error())
		response.FailResponseWithMsgErr("服务器错误", err, c)
		return
	}
	if code != emailCode {
		response.SuccessResponseWithMsg("验证码错误", c)
		return
	}
	//  创建用户
	u := models.UserBasic{
		Identity: utils.GetUUID(),
		UserName: username,
		Password: utils.GetMd5(password),
		Email:    email,
	}
	fmt.Println(u)
	err = utils.DB.Create(&u).Error
	if err != nil {
		response.FailResponseWithMsg("系统错误", c)
		return
	}
	fmt.Println(u)
	response.SuccessResponse("注册成功", u, c)
}

// SendEmailCode
// @Tags 公共方法
// @Summary 用户注册
// @Param email formData string false "email"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /email/code [post]
func SendEmailCode(c *gin.Context) {
	//获取邮箱
	email := c.PostForm("email")
	if email == "" {
		response.FailResponseWithMsg("邮箱为空", c)
		return
	}
	//校验邮箱格式
	if ok := utils.IsEmailValid(email); !ok {
		response.FailResponseWithMsg("请输入正确的邮箱", c)
		return
	}
	u := models.UserBasic{}
	//查看邮箱是否被注册
	count := utils.DB.Where("email = ?", email).Find(&u).RowsAffected
	if count > 0 {
		response.FailResponseWithMsg("邮箱已被注册", c)
		return
	}

	//fmt.Println(u)
	//生成4位数验证码
	code := utils.GetCode()
	fmt.Println(code)
	//存到redis
	err := utils.Redis.Set(context.Background(), global.EmailCode+email, code, time.Second*60).Err()
	if err != nil {
		log.Printf("%v", err.Error())
		err = utils.SendEmailCode(email, code)
		if err != nil {
			response.FailResponseWithMsgErr("服务器错误", err, c)
			response.FailResponseWithMsgErr("服务器错误", err, c)
			return
		}
		//发送验证码
		return
	}
	response.SuccessResponseWithMsg("已发送验证码，一分钟内有效，请注意查收", c)
}

// UserDetail
// @Tags 用户组
// @Summary 用户信息详情
// @Success 200 {string} json "{"code":"200","msg":"",data:""}"
// @Router /u/detail [get]
func UserDetail(c *gin.Context) {
	//在auth中间件中 保存了用户认证信息 user_claim
	u, _ := c.Get("user_claim")
	//断言
	uc := u.(*utils.UserClaim)

	//	查数据库
	userDetail := &models.UserBasic{}
	tx := models.GetUserBasicDetail(uc.Identity)
	err := tx.Find(&userDetail).Error
	if err != nil {
		response.FailResponseWithMsg("服务器错误", c)
		return
	}
	//	返回结果
	//fmt.Println(problems)
	response.SuccessResponseWithData(gin.H{"data": userDetail}, c)
}
