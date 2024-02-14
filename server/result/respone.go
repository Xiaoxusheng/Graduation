package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ok 成功
func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  "success"})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg})
}
