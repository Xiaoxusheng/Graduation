package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"time"
)

// Log 记录日志中间件
func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	获取请求信息
		path := c.Request.URL.Path
		method := c.Request.Method

		t := time.Now()
		c.Next()
		httpCode := c.Writer.Status()
		id := c.GetString("identity")
		//	写入数据库
		global.Global.Log.Info(time.Now().Sub(t))
		fmt.Println(path, method, id, httpCode)
		if path == "" || method == "" || id == "" || httpCode == 0 {
			return
		}
		fmt.Println(path, method, id, httpCode)

	}
}
