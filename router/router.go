package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"online_exercise_system/middleware"
	"online_exercise_system/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/hello", service.Hello)
	r.POST("/login", service.Login)
	r.POST("/email/code", service.SendEmailCode)
	r.POST("/register", service.Register)
	//用户提交排行
	r.GET("/rank_list", service.GetRankList)
	//用户提交列表
	r.GET("/submit_list", service.SubmitList)

	problem := r.Group("/problem")
	//问题列表
	problem.GET("/list", service.ProblemList)
	//问题详情
	problem.GET("/detail", service.ProblemDetail)

	//用户组
	authUser := r.Group("/u", middleware.AuthCheck())
	authUser.GET("/detail", service.UserDetail)
	//提交问题
	authUser.POST("/submit", service.SubmitCode)

	authAdmin := authUser.Group("/admin", middleware.AuthAdminCheck())
	//管理员创建问题
	authAdmin.POST("/problem_add", service.CreateProblem)
	//管理员修改问题
	authAdmin.POST("/problem_modify", service.ModifyProblem)
	//管理员获取分类列表
	authAdmin.GET("/category_list", service.CategoryList)
	//管理员添加分类
	authAdmin.POST("/category_add", service.CategoryAdd)
	//管理员删除分类
	authAdmin.DELETE("/category_del", service.CategoryDel)
	//管理员更新分类
	authAdmin.PUT("/category_update", service.CategoryUpdate)

	authUser.POST("/test", service.Hello)

	return r
}
