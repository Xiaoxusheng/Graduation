package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strconv"
	"time"
)

/*考勤模块*/

// ClockIn 打卡
func ClockIn(c *gin.Context) {
	//考虑时间问题
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	var status int32 = 2
	t1 := time.Now()
	//考虑重复打卡
	if global.Global.Redis.SIsMember(global.Global.Ctx, global.ClockIn, id).Val() {
		//  已经打卡一次
		//获取uid
		val := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
		if val != "" {
			n, err := strconv.Atoi(val)
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.ClockInError)
				return
			}
			//打卡多晚只有5.20，加班另外算
			err = dao.AfterWork(int64(n), t1)
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.ClockInError)
				return
			}
			result.Ok(c, nil)
			err = global.Global.Pool.Submit(func() {
				//删除之前的记录
				_, err = global.Global.Redis.Del(global.Global.Ctx, global.GetClockInLog+val).Result()
				if err != nil {
					global.Global.Log.Error(err)
				}
			})
			return
		}
		//数据库查
		employer, err := dao.GetUserById(id)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		//插入
		err = dao.AfterWork(employer.Uid, t1)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		err = global.Global.Pool.Submit(func() {
			//插入uid与identity的关系
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, employer.Identity, employer.Uid).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
			//设置过期时间
			_, err = global.Global.Redis.Expire(global.Global.Ctx, global.UidId, time.Second*global.EmployerUidId).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
			//删除之前的记录
			_, err = global.Global.Redis.Del(global.Global.Ctx, global.GetClockInLog+strconv.FormatInt(employer.Uid, 10)).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
		})
		if err != nil {
			global.Global.Log.Error(err)
		}
		result.Ok(c, nil)
		return
	}
	//第一次打卡
	//获取uid
	val := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	v4 := utils.GetUidV4()

	if val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		if t1.Unix() > time.Date(t1.Year(), t1.Month(), t1.Day(), 9, 0, 0, 0, t1.Location()).Unix() {
			status = 3
		}
		err = dao.ClockIn(&models.Attendance{
			Identity:  v4,
			Uid:       int64(n),
			Date:      t1,
			StartTime: t1,
			Status:    status,
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		err = global.Global.Pool.Submit(func() {
			global.Global.Wg.Add(1)
			defer global.Global.Wg.Done()
			//打卡
			global.Global.Redis.SAdd(global.Global.Ctx, global.ClockIn, id)
			_, err = global.Global.Redis.Expire(global.Global.Ctx, global.ClockIn, time.Date(t1.Year(), t1.Month(), t1.Day()+1, 0, 0, 0, 0, t1.Location()).Sub(t1)).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
			_, err = global.Global.Redis.Del(global.Global.Ctx, global.GetClockInLog+val).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
		})
		if err != nil {
			global.Global.Log.Error(err)
		}
		result.Ok(c, nil)
		return
	}
	//数据库查
	employer, err := dao.GetUserById(id)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.ClockInError)
		return
	}
	//插入
	if t1.Unix() > time.Date(t1.Year(), t1.Month(), t1.Day(), 9, 0, 0, 0, t1.Location()).Unix() {
		status = 3
	}
	err = dao.ClockIn(&models.Attendance{
		Identity:  v4,
		Uid:       employer.Uid,
		Date:      t1,
		StartTime: t1,
		Status:    status,
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.ClockInError)
		return
	}

	go func() {
		//插入uid与identity的关系
		_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, employer.Identity, employer.Uid).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
		//设置过期时间
		_, err = global.Global.Redis.Expire(global.Global.Ctx, global.UidId, time.Second*global.EmployerUidId).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
		//打卡名单
		_, err = global.Global.Redis.SAdd(global.Global.Ctx, global.ClockIn, id).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
		//第二天0点删除打卡
		_, err = global.Global.Redis.Expire(global.Global.Ctx, global.ClockIn, time.Date(t1.Year(), t1.Month(), t1.Day()+1, 0, 0, 0, 0, t1.Location()).Sub(t1)).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
		_, err = global.Global.Redis.Del(global.Global.Ctx, global.GetClockInLog+val).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}

	}()
	result.Ok(c, nil)
}

// MarkCardApplication 补卡申请
func MarkCardApplication(c *gin.Context) {
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
		//处理
		markCard := new(global.MarkCard)
		err = c.Bind(markCard)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataConflict, global.QueryNotFoundError)
			return
		}
		//判断员工是否在职
		if employer.Status != 1 {
			result.Fail(c, global.DataConflict, global.UserNotWorkError)
			return
		}
		//插入数据
		err = dao.InsertMarkCardApplication(&models.Examine{
			Identity: utils.GetUidV4(),
			Uid:      employer.Uid,
			Status:   1,
			Reason:   markCard.Reason,
			Url:      markCard.Url,
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.MarkCardError)
			return
		}
		result.Ok(c, nil)
		//	同步到redis
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
		return
	}
	//缓存中存在
	uids, err := strconv.Atoi(uid)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	markCard := new(global.MarkCard)
	err = c.Bind(markCard)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataUnmarshal, global.QueryNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//插入数据
	err = dao.InsertMarkCardApplication(&models.Examine{
		StartTime: time.Unix(markCard.Date, 0).Add(time.Hour * 9),
		EndTime:   time.Unix(markCard.Date, 0).Add(time.Hour * 17).Add(time.Minute * 20),
		Identity:  utils.GetUidV4(),
		Uid:       int64(uids),
		Status:    2,
		Reason:    markCard.Reason,
		Url:       markCard.Url,
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.MarkCardError)
		return
	}
	result.Ok(c, nil)
}

// LeaveApplication 请假申请
func LeaveApplication(c *gin.Context) {
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
			result.Fail(c, global.ServerError, global.LeaveApplicationError)
			return
		}
		leave := new(global.LeaveApplication)
		err = c.Bind(leave)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.QueryNotFoundError)
			return
		}
		//判断员工是否在职
		if employer.Status != 1 {
			result.Fail(c, global.DataConflict, global.UserNotWorkError)
			return
		}

		err = dao.InsertMarkCardApplication(&models.Examine{
			Identity:  utils.GetUidV4(),
			Uid:       employer.Uid,
			StartTime: time.Unix(leave.StartTime, 0),
			EndTime:   time.Unix(leave.EndTime, 0),
			Status:    4,
			Url:       leave.Url,
			Reason:    leave.Reason,
			Model:     gorm.Model{},
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.LeaveError)
			return
		}
		//	同步到redis
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
		result.Ok(c, nil)
		return
	}
	uids, err := strconv.Atoi(uid)
	if err != nil {
		result.Fail(c, global.ServerError, global.AtoiError)
		global.Global.Log.Error(err)
		return
	}
	leave := new(global.LeaveApplication)
	err = c.Bind(leave)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataUnmarshal, global.QueryNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	err = dao.InsertMarkCardApplication(&models.Examine{
		Identity:  utils.GetUidV4(),
		Uid:       int64(uids),
		Url:       leave.Url,
		StartTime: time.Unix(leave.StartTime, 0),
		EndTime:   time.Unix(leave.EndTime, 0),
		Status:    4,
		Reason:    leave.Reason,
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataUnmarshal, global.LeaveError)
		return
	}
	result.Ok(c, nil)
}

// OverTimeApplication 加班申请
func OverTimeApplication(c *gin.Context) {
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
			result.Fail(c, global.ServerError, global.OverTimeError)
			return
		}
		overtime := new(global.OverTime)
		err = c.Bind(overtime)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.QueryNotFoundError)
			return
		}
		//判断员工是否在职
		if employer.Status != 1 {
			result.Fail(c, global.DataConflict, global.UserNotWorkError)
			return
		}

		err = dao.InsertMarkCardApplication(&models.Examine{
			Identity:  utils.GetUidV4(),
			Uid:       employer.Uid,
			Url:       overtime.Url,
			StartTime: time.Unix(overtime.StartTime, 0),
			EndTime:   time.Unix(overtime.EndTime, 0),
			Status:    1,
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.LeaveError)
		}
		//	同步到redis
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
			result.Fail(c, global.DataUnmarshal, global.OverTimeError)
			return
		}
		result.Ok(c, nil)
	}
	uids, err := strconv.Atoi(uid)
	if err != nil {
		result.Fail(c, global.ServerError, global.AtoiError)
		global.Global.Log.Error(err)
		return
	}
	overtime := new(global.OverTime)
	err = c.Bind(overtime)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataUnmarshal, global.QueryNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//插入数据库
	err = dao.InsertMarkCardApplication(&models.Examine{
		Identity:  utils.GetUidV4(),
		Uid:       int64(uids),
		Url:       overtime.Url,
		StartTime: time.Unix(overtime.StartTime, 0),
		EndTime:   time.Unix(overtime.EndTime, 0),
		Status:    1,
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataUnmarshal, global.OverTimeError)
		return
	}
	result.Ok(c, nil)
}

func GetClockInLog(c *gin.Context) {
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	var uid string
	uid = global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	if uid == "" {
		employer, err := dao.GetUserById(id)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.OverTimeError)
			return
		}
		//判断员工是否在职
		if employer.Status != 1 {
			result.Fail(c, global.DataConflict, global.UserNotWorkError)
			return
		}
		uid = strconv.FormatInt(employer.Uid, 10)
		//	同步到redis
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
			result.Fail(c, global.ServerError, global.GetClockError)
			return
		}
		return
	}
	//获取考勤列表
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")
	uids, err := strconv.Atoi(uid)
	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		result.Fail(c, global.DataUnmarshal, global.AtoiError)
		global.Global.Log.Error(err)
		return
	}
	//查询uid是否存在 工号
	if !global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uids).Val() {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//先在缓存中获取,存的是员工考勤记录
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
	list, err := dao.GetAttendanceList(limits, offsets, int64(uids))
	if err != nil {
		result.Fail(c, global.ServerError, global.GetClockError)
		global.Global.Log.Error(err)
		return
	}
	//放进redis
	err = global.Global.Pool.Submit(func() {
		marshal, err := json.Marshal(list)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.GetClockInLog+uid, marshal, time.Second*global.EmployerClockTime).Result()
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetClockError)
		return
	}
	result.Ok(c, list)
}
