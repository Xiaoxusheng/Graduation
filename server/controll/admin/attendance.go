package admin

import (
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/result"
	"strconv"
)

/*管理端*/

// GetClockIn  获取获取考勤记录
func GetClockIn(c *gin.Context) {
	//工号
	uid := c.Query("id")
	if uid == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	uids, err := strconv.Atoi(uid)
	if err != nil {
		global.Global.Log.Error(err)
		return
	}
	//查询id是否存在 工号
	list, err := dao.GetAttendanceList(int32(uids))
	if err != nil {
		global.Global.Log.Error(err)
		return
	}
	result.Ok(c, list)
}

//编辑考勤记录
