package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/config"
	"server/middleware"
	"server/router"
	"strconv"
	"time"
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
	//初始化锁
	config.InitMutex()

	e := gin.Default()

	//中间件
	e.Use(middleware.Log(), middleware.Cors(), middleware.RateLimit())

	e = router.Routers(e)

	s := &http.Server{
		Addr:           ":" + strconv.Itoa(config.Config.Service.Port),
		Handler:        e,
		ReadTimeout:    config.Config.Service.ReadTime * time.Second,
		WriteTimeout:   config.Config.Service.WriteTime * time.Second,
		MaxHeaderBytes: int(config.Config.Service.MaxHeaderBytes << 20),
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Println("server start fail")
	}

	//err := e.Run(":80")
	//if err != nil {
	//	log.Println("server start fail")
	//	return
	//}
}
