package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/result"
	"strconv"
	"time"
)

// GetExamine 员工查看申请状态
func GetExamine(c *gin.Context) {
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
	//缓存中查
	val := global.Global.Redis.Get(global.Global.Ctx, global.Examine+uid).Val()
	if val != "" {
		list := make([]*global.Applications, 0)
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.DataUnmarshalError)
			return
		}
		result.Ok(c, list)
		return
	}
	list, err := dao.GetExamine(uid)
	if err != nil {
		return
	}
	err = global.Global.Pool.Submit(func() {
		marshal, err := json.Marshal(list)
		if err != nil {
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Examine+uid, marshal, global.ExamineTime*time.Second).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, list)
}
