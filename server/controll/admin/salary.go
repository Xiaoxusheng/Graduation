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
	if u == "" {
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
	//获取
	val := global.Global.Redis.Get(global.Global.Ctx, global.SalaryEmployerList+u).Val()
	if val != "" {
		list := make([]*models.Salary, 0)
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.DataUnmarshalError)
			return
		}
		result.Ok(c, list)
		return
	}
	salaryList, err := dao.GetSalaryByEmployer(int64(uid))
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetSalaryError)
		return
	}
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		//	同步
		marshal, err := json.Marshal(salaryList)
		if err != nil {
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.SalaryEmployerList+u, marshal, global.SalaryEmployerListTime*time.Second).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, salaryList)
}

//定时任务取计算工资

// GetSalaryList 获取某个月工资表
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
	val := global.Global.Redis.HGet(global.Global.Ctx, global.SalaryList, t).Val()
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
		marshal, err := json.Marshal(list)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.HSet(global.Global.Ctx, global.SalaryList, t, marshal).Result()
		global.Global.Redis.Expire(global.Global.Ctx, global.SalaryList, global.SalaryListTime*time.Second)
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, list)
}

// SalaryInfo 输入员工时长
func SalaryInfo(c *gin.Context) {
	salaryInfo := new(global.SalaryInfos)
	err := c.Bind(salaryInfo)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, salaryInfo.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//	存入数据库
	err = dao.InsertSalary(&models.Salary{
		Identity:        utils.GetUidV4(),
		Uid:             salaryInfo.Uid,
		PaidLeave:       salaryInfo.PaidLeave,
		Count:           salaryInfo.Count,
		AttendanceHours: salaryInfo.AttendanceHours,
		Total:           float64(salaryInfo.PaidLeave*20) + salaryInfo.AttendanceHours*20 + salaryInfo.Other - salaryInfo.Other,
		Other:           salaryInfo.Other,
		Subsidy:         salaryInfo.Subsidy,
		Date:            time.Now().Format("2006-01"),
	})

	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.InputSalaryError)
		return
	}
	result.Ok(c, nil)
}

// DeleteSalary  删除员工工资信息
func DeleteSalary(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//id不存在
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.SalaryId, id).Val() {
		result.Fail(c, global.DataConflict, global.DelSalaryError)
		return
	}
	err := dao.DeleteSalary(id)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.DelSalaryError)
		return
	}
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		_, err = global.Global.Redis.Del(global.Global.Ctx, global.SalaryList).Result()
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.SRem(global.Global.Ctx, global.SalaryId, id).Result()
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
	})
	if err != nil {
		global.Global.Log.Error("goroutine fail:", err)
	}
	result.Ok(c, nil)
}
