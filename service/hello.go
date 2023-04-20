package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /hello [get]
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "hello", "code": 200})
}
