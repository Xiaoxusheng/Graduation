package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
	"strconv"
	"strings"
	"time"
)

/*考勤模块*/

// ClockIn 打卡
func ClockIn(c *gin.Context) {
	//考虑重复打卡

	//考虑时间问题
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	d := time.Now().Day()
	global.Global.Log.Warn(d)
	val := global.Global.Redis.BitField(global.Global.Ctx, global.Sign+id, "GET", "u"+strconv.Itoa(d), 0).Val()
	global.Global.Log.Info(val)
	s := strings.Builder{}
	if val[0] == 1 {
		for i := 1; i < d; i++ {
			s.WriteString("0")
		}
		s.WriteString("1")
		result.Ok(c, s.String())
		return
	}
	result.Ok(c, fmt.Sprintf("%b", val[0]))
	return
}

//补卡
