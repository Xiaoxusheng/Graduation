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

// AddEmployee 添加员工信息
func AddEmployee(c *gin.Context) {
	e := new(global.Employers)
	e.Uid = 123
	err := c.Bind(e)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}

	//生成唯一员工标识
	uid, err := utils.GetUid(e.Department)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}

	if err != nil {
		global.Global.Log.Error(err)
		return
	}
	//	添加员工信息
	id := utils.GetUidV4()
	err = dao.InsertEmployer(&models.Employee{
		Identity:     id,
		Uid:          uid,
		Name:         e.Name,
		Birthday:     time.Unix(e.Birthday, 0),
		Phone:        e.Phone,
		Status:       1,
		Position:     e.Position,
		DepartmentId: e.Department,
		IP:           c.RemoteIP(),
	}, func() {
		val := global.Global.Redis.HGet(global.Global.Ctx, global.Uid, string(e.Department)).Val()
		if val != "" {
			num, err := strconv.Atoi(val)
			if err != nil {
				global.Global.Log.Error(err)
			}
			num = num - 1
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.Uid, string(e.Department), strconv.Itoa(num)).Result()
			if err != nil {
				global.Global.Log.Error(err)
			}
		}
	})
	//global.Global.Redis.Set(global.Global.Ctx, global.UidKey+use.Identity, use., 0)

	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.AddEmployerError)
		return
	}

	//添加进缓存
	go func() {
		//uid添加进缓存
		_, err = global.Global.Redis.SAdd(global.Global.Ctx, global.Employer, uid).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()

	result.Ok(c, nil)
}

// DeleteEmployee 删除员工信息
func DeleteEmployee(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	e, err := dao.GetEmployerByUid(uid)
	if e == nil || err != nil {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//	删除员工信息
	err = dao.DeleteEmployer(uid)
	if err != nil {
		result.Fail(c, global.BadRequest, global.DeleteError)
		return
	}
	err = dao.DeleteUser(e.Identity)
	if err != nil {
		result.Fail(c, global.BadRequest, global.DeleteError)
		return
	}
	go func() {
		_, err = global.Global.Redis.SAdd(global.Global.Ctx, global.Employer, e.Uid).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	result.Ok(c, nil)
}

// EmployeeList 员工列表
func EmployeeList(c *gin.Context) {
	//分页查询
	offset := c.DefaultQuery("offset", "1")
	limit := c.DefaultQuery("limit", "10")
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	limits, err := strconv.Atoi(limit)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	list, err := dao.GetEmployerList(limits, offsets)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	result.Ok(c, list)
}

// UpdateEmployee 更新员工信息
func UpdateEmployee(c *gin.Context) {
	e := new(global.Employers)
	err := c.Bind(e)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryError)
		global.Global.Log.Error(err)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, e.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//	更新信息
	err = dao.UpdateEmployer(e)
	if err != nil {
		result.Fail(c, global.ServerError, global.UpdateEmployerError)
		global.Global.Log.Error(err)
		return
	}
	go func() {
		_, err = global.Global.Redis.Del(global.Global.Ctx, global.Uid+strconv.Itoa(int(e.Uid))).Result()
		if err != nil {
			global.Global.Log.Error(e)
		}
	}()

	result.Ok(c, nil)
}

// EmployeeInfo 获取员工信息
func EmployeeInfo(c *gin.Context) {
	uid, err := strconv.Atoi(c.Query("uid"))
	if uid == 0 || err != nil {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//获取员工信息
	val := global.Global.Redis.Get(global.Global.Ctx, global.Uid+strconv.Itoa(uid)).Val()
	if val != "" {
		e := new(models.Employee)
		err = json.Unmarshal([]byte(val), e)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		result.Ok(c, e)
		return
	}
	//不存在，过期
	info, err := dao.GetEmployerInfo(int64(uid))
	if err != nil {
		result.Fail(c, global.DataNotFound, global.UserNotExistError)
		return
	}
	//插入
	go func() {
		if info == nil {
			_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+strconv.Itoa(uid), "null", global.InfoTime).Result()
			return
		}
		marshal, err := json.Marshal(info)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+strconv.Itoa(uid), marshal, 0).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	result.Ok(c, info)
}
