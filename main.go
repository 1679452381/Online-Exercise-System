package main

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "online_exercise_system/docs"
	"online_exercise_system/router"
)

// @title online_exercise_system
// @version 1.0
// @description  API文档
// @host   127.0.0.1:8080
// @BasePath /
func main() {
	r := router.Router()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
