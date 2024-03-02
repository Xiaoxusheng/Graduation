package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strconv"
	"time"
)

// GetSalary 查员工的工资
func GetSalary(c *gin.Context) {
	u := c.Query("uid")
	date := c.Query("date")
	if date == "" || u == "" {
		global.Global.Log.Warn("uid is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	uid, err := strconv.Atoi(u)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	//	判断员工是否存在
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val() {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	t, err := time.Parse("2006-01", date)
	if err != nil {
		result.Fail(c, global.ServerError, err.Error())
		global.Global.Log.Error(err)
		return
	}
	//时间比当前时间大
	if t.After(time.Now()) {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	var times float64
	var count int32
	salary, err := dao.GetSalary(int64(uid), date)
	if err != nil {
		//去获取数据插入salary表
		list, err := dao.GetAttendance(int64(uid), t)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.GetSalaryError)
			return
		}

		for i := 0; i < len(list); i++ {
			times = times + list[i].EndTime.Sub(list[i].StartTime).Hours()
			if list[i].EndTime.Hour() < 17 && list[i].EndTime.Month() < 20 {
				//早退
				count++
			}
			if list[i].StartTime.Hour() > 9 && list[i].StartTime.Month() > 0 {
				//迟到
				count++
			}
		}
		// 获取当前时间
		now := time.Now()
		// 获取当前月份的天数
		//year, month, _ := now.Date()
		daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		h := time.Date(now.Year(), now.Month(), now.Day(), 17, 20, 0, 0, time.Local).Sub(time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)).Hours()
		model := &models.Salary{
			Identity:                utils.GetUidV4(),
			Uid:                     int64(uid),
			PaidLeave:               0,
			Count:                   count,
			Total:                   0,
			ExpectedAttendanceHours: float64(daysInMonth) * h,
			AttendanceHours:         times,
			Other:                   0,
			Subsidy:                 0,
			Date:                    t.Format("2006-01"),
		}
		err = dao.InsertSalary(model)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.GetSalaryError)
			return
		}
		result.Ok(c, model)
		return
	}
	//超过
	if time.Now().Sub(salary.CreatedAt).Hours() > 24 {
		//	更新
		//去获取数据插入salary表
		list, err := dao.GetAttendance(int64(uid), t)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.GetSalaryError)
			return
		}
		for i := 0; i < len(list); i++ {
			times = times + list[i].EndTime.Sub(list[i].StartTime).Hours()
			if list[i].EndTime.Hour() < 17 && list[i].EndTime.Month() < 20 {
				//早退
				count++
			}
			if list[i].StartTime.Hour() > 9 && list[i].StartTime.Month() > 0 {
				//迟到
				count++
			}
		}
		err = dao.UpdateSalary(int64(uid), count, times)
		if err != nil {
			return
		}
		list, err = dao.GetAttendance(int64(uid), t)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.GetSalaryError)
			return
		}
		result.Ok(c, list)
		return
	}

	result.Ok(c, salary)

}

//定时任务取计算工资

// GetSalaryList 获取员工工资表
func GetSalaryList(c *gin.Context) {
	t := c.DefaultQuery("time", time.Now().Format("2006-01"))
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")

	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	val := global.Global.Redis.Get(global.Global.Ctx, global.SalaryList).Val()
	if val != "" {
		list := make([]*models.Salary, 0)
		err = json.Unmarshal([]byte(val), &list)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.DataUnmarshalError)
			return
		}
		result.Ok(c, list)
		return
	}
	list, err := dao.GetSalaryList(limits, offsets, t)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetSalaryListError)
		return
	}
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		marshal, err := json.Marshal(limit)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.SalaryList, marshal, global.SalaryListTime).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, list)
}
