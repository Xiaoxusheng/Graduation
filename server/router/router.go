package router

import (
	"github.com/gin-gonic/gin"
	"server/controll/user"
	"server/middleware"
)

func Routers(e *gin.Engine) *gin.Engine {
	//

	//登录
	e.POST("/user/login", user.Login)
	//注册
	e.POST("/user/register", user.Register)
	api := e.Group("/user")
	api.Use(middleware.ParseToken())
	//个人信息
	api.GET("/info", user.Info)
	//退出登录
	api.GET("/Logout", user.Logout)

	//删除员工
	api.GET("delete_employer", user.DeleteEmployee)
	//创建员工信息
	api.POST("/add_employer", user.AddEmployee)
	//更新员工信息

	//员工信息列表

	//员工个人信息
	api.GET("info", user.EmployeeInfo)
	return e
}
