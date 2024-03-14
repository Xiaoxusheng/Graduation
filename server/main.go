package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/config"
	"server/global"
	"server/middleware"
	"server/router"
	"server/utils"
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
	//初始化协程池
	config.InitPool()
	//初始化casBin
	config.InitCasBin()

	e := gin.Default()

	//中间件
	e.Use(middleware.Cors(), middleware.RateLimit())

	e = router.Routers(e)
	go utils.Listen()
	global.Global.Wg.Wait()

	s := &http.Server{
		Addr:           ":" + strconv.Itoa(config.Config.Service.Port), //端口
		Handler:        e,                                              //路由
		ReadTimeout:    config.Config.Service.ReadTime * time.Second,   //读超时时间
		WriteTimeout:   config.Config.Service.WriteTime * time.Second,  //写超时时间
		MaxHeaderBytes: config.Config.Service.MaxHeaderBytes << 20,     //最大请求头 1MB
	}
	err := s.ListenAndServe()
	//err := s.ListenAndServeTLS("/root/ssl/xlei.love.pem", "/root/ssl/xlei.love.key")
	if err != nil {
		log.Println("server start fail :" + err.Error())
	}

	//err := e.Run(":80")
	//if err != nil {
	//	log.Println("server start fail")
	//	return
	//}
}
