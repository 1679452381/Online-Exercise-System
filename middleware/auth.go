package middleware

import (
	"github.com/gin-gonic/gin"
	"online_exercise_system/response"
	"online_exercise_system/utils"
)

// 用户登录验证
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	获取token
		token := c.GetHeader("authorization")
		//	解析token
		uc, err := utils.AnalyToken(token)
		if err != nil {
			c.Abort()
			response.FailResponseUnauthorizedWithMsg("用户认证失败", c)
			return
		}
		//从redis获取token
		//rToken, err := utils.Redis.Get(context.Background(), global.Token+uc.Identity).Result()
		//if err != nil {
		//	c.Abort()
		//	response.FailResponseWithMsg("登陆超时，请重新登录", c)
		//	return
		//}
		//if token != rToken {
		//	c.Abort()
		//	response.FailResponseWithMsg("登陆超时，请重新登录", c)
		//	return
		//}
		c.Set("user_claim", uc)
		c.Next()
	}
}

// 管理员登录验证
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := c.Get("user_claim")
		uc := u.(*utils.UserClaim)
		if uc.IsAdmin != 1 {
			c.Abort()
			response.FailResponseUnauthorizedWithMsg("非管理员用户", c)
			return
		}
		c.Next()
	}
}
