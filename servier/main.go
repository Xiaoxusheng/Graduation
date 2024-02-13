package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//初始化配置

	e := gin.Default()

	//中间件

	err := e.Run(":80")
	if err != nil {
		log.Println("server start fail")
		return
	}
}
