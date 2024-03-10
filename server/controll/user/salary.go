package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"strconv"
)

func GetSalary(c *gin.Context) {
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	uid := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	if uid == "" {
		//数据库查
		employer, err := dao.GetUserById(id)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		uid = strconv.FormatInt(employer.Uid, 10)
		//	插入redis
		err = global.Global.Pool.Submit(func() {
			global.Global.Wg.Add(1)
			defer global.Global.Wg.Done()
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, id, employer.Uid).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
		})
		if err != nil {
			global.Global.Log.Error("submit err :", err)
		}
	}
	//	判断员工是否存在
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val() {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//获取
	val := global.Global.Redis.Get(global.Global.Ctx, global.SalaryEmployerList+uid).Val()
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
	uids, err := strconv.Atoi(uid)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	salaryList, err := dao.GetSalaryByEmployer(int64(uids))
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
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.SalaryEmployerList+uid, marshal, global.SalaryEmployerListTime).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, salaryList)
}
