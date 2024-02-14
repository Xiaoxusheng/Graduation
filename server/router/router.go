package router

import (
	"github.com/gin-gonic/gin"
	"server/controll/user"
)

func Router(e *gin.Engine) *gin.Engine {
	//
	api := e.Group("/user")
	api.POST("/login", user.Login)
	api.POST("/register", user.Register)
	api.GET("/info", user.Info)

	return e
}
