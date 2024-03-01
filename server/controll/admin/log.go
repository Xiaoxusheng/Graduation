package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/dao"
	"server/global"
	"server/result"
	"strconv"
)

func GetLogList(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")

	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	var count = new(int64)
	list, err := dao.GetLogList(limits, offsets)
	err = dao.GetCount(count)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetLogList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  list,
		"count": count,
		"msg":   "success",
	})
	//result.Ok(c, list)
}
