package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
)

func CasBin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetString("identity")
		path := c.FullPath()
		method := c.Request.Method
		fmt.Println(user, path, method)
		ok, err := global.Global.CasBin.Enforce(user, path, method)
		if err != nil || !ok {
			global.Global.Log.Error("权限验证失败：", err)
			result.Fail(c, global.PermissionDenied, global.PermissionDeniedError)
			c.Abort()
			return
		}
		c.Next()
	}
}
