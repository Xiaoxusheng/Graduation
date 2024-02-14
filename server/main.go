package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"server/config"
	"server/router"
)

func main() {
	//日志颜色开启
	gin.ForceConsoleColor()
	//初始化配置
	config.InitService()
	// 记录到文件。
	config.InitLog()
	//初始化mysql
	config.InitMysql()
	//初始化redis
	config.InitRedis()

	e := gin.Default()

	//中间件

	e = router.Router(e)

	err := e.Run(":80")
	if err != nil {
		log.Println("server start fail")
		return
	}
}
