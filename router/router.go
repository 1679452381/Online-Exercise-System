package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"online_exercise_system/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/hello", service.Hello)
	r.POST("/login", service.Login)
	return r
}
