package result

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// Ok 成功
func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  "success"})
}

// Fail 失败
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg})
}

// Stream 流式响应
func Stream(c *gin.Context, reader *bufio.Reader) {
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, reader)
		if err != nil {
			return false
		}
		return true
	})
}
