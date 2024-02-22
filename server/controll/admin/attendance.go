package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"strconv"
	"time"
)

/*考勤模块*/

// GetClockInLog  获取员工的获取考勤记录
func GetClockInLog(c *gin.Context) {
	//工号
	uid := c.Query("uid")
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
	//查询uid是否存在 工号
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uids).Val() {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//先在缓存中获取,存的是员工工号
	val := global.Global.Redis.Get(global.Global.Ctx, global.GetClockInLog+uid).Val()
	if val != "" {
		//	val存在
		list := make([]*models.Attendance, 0)
		err = json.Unmarshal([]byte(val), &list)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.DataUnmarshalError)
			return
		}
		result.Ok(c, list)
		return
	}
	//数据库中获取
	list, err := dao.GetAttendanceList(int32(uids))
	if err != nil {
		global.Global.Log.Error(err)
		return
	}
	go func() {
		marshal, err := json.Marshal(list)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		global.Global.Redis.Set(global.Global.Ctx, global.GetClockInLog+uid, marshal, global.EmployerInfoTime)
	}()

	result.Ok(c, list)

}

// EditClockLog 编辑考勤记录
func EditClockLog(c *gin.Context) {
	attendances := new(global.Attendance)
	err := c.Bind(attendances)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.ParseError)
		return
	}
	//验证数据
	num := map[int32]bool{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
	}
	if _, ok := num[attendances.Status]; !ok {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//	修改
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()
	err = dao.UpdateAttendance(attendances)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.UpdateClockInLogError)
		return
	}
	result.Ok(c, nil)
}

// GetClockList 获取某一天员工所有员工的打卡记录
func GetClockList(c *gin.Context) {
	//	获取当前时间
	//
	t := c.Query("time")
	t1, err := strconv.Atoi(t)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	date := time.Unix(int64(t1), 0)
	if t1 < 0 {
		result.Fail(c, global.DataValidationError, global.QueryError)
		return
	}
	list, err := dao.GetDateList(date)
	if err != nil {
		result.Fail(c, global.ServerError, global.GetClockError)
		return
	}
	result.Ok(c, list)
}

// ClockInSummaryList 考勤总结审核
func ClockInSummaryList(c *gin.Context) {

}

//
