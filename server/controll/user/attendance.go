package user

import (
	"github.com/gin-gonic/gin"
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
		val := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
		if val != "" {
			n, err := strconv.Atoi(val)
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.ClockInError)
				return
			}
			//打卡多晚只有5.20，加班另外算
			err = dao.AfterWork(int64(n), time.Date(t1.Year(), t1.Month(), t1.Day(), 17, 20, 0, 0, t1.Location()).Unix())
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.ClockInError)
				return
			}
			global.Global.Redis.SAdd(global.Global.Ctx, global.ClockIn, id)
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
		err = dao.AfterWork(employer.Uid, t1.Unix())
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
		}()
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
			StartTime: t1,
			Status:    status,
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		global.Global.Redis.SAdd(global.Global.Ctx, global.ClockIn, id)
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
		_, err = global.Global.Redis.Expire(global.Global.Ctx, global.ClockIn, time.Duration(time.Date(t1.Year(), t1.Month(), t1.Day()+1, 0, 0, 0, 0, t1.Location()).Sub(t1).Seconds())).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}

	}()
	result.Ok(c, nil)
	//d := time.Now().Day()
	//global.Global.Log.Warn(d)
	//val := global.Global.Redis.BitField(global.Global.Ctx, global.Sign+id, "GET", "u"+strconv.Itoa(d), 0).Val()
	//global.Global.Log.Info(val)
	//s := strings.Builder{}
	//if val[0] == 1 {
	//	for i := 1; i < d; i++ {
	//		s.WriteString("0")
	//	}
	//	s.WriteString("1")
	//	result.Ok(c, s.String())
	//	return
	//}
	//result.Ok(c, fmt.Sprintf("%b", val[0]))
	//return
}

// MarkCardApplication 补卡申请
func MarkCardApplication(c *gin.Context) {
	markCard := new(global.MarkCard)
	err := c.Bind(markCard)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, markCard.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//插入数据
	err = dao.InsertMarkCardApplication(&models.Examine{
		Identity: utils.GetUidV4(),
		Name:     "",
		Uid:      markCard.Uid,
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
}

//请假

//加班申请

//
