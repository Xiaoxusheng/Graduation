package router

import (
	"github.com/gin-gonic/gin"
	"server/controll/user"
)

func Routers(e *gin.Engine) *gin.Engine {
	//
	api := e.Group("/user")
	//登录
	api.POST("/login", user.Login)
	//注册
	api.POST("/register", user.Register)
	//个人信息
	api.GET("/info", user.Info)
	//退出登录
	api.GET("/Logout", user.Logout)

	return e
}
