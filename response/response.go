package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessResponseWithMsg(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 200,
	})
}

func SuccessResponseWithData(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  "请求成功",
		"code": 200,
		"data": data,
	})
}
func SuccessResponse(msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 200,
		"data": data,
	})
}
func SuccessResponseWithToken(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":   "登陆成功",
		"code":  200,
		"token": data,
	})
}
func FailResponseWithMsg(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": -1,
	})
}
