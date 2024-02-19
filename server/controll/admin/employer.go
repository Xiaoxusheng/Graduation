package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strconv"
	"time"
)

type employer struct {
	Name       string `json:"name,omitempty"  binding:"required" form:"name"`
	Birthday   int64  `json:"birthday,omitempty" binding:"required" form:"birthday"`
	Phone      string `json:"phone,omitempty" binding:"required" form:"phone"`
	Position   int32  `json:"position,omitempty" binding:"required,number" form:"position"`
	Department int32  `json:"department,omitempty" binding:"required,number" form:"department"`
}

// AddEmployee 添加员工信息
func AddEmployee(c *gin.Context) {
	e := new(employer)
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
	fmt.Println(time.Unix(e.Birthday, 0))
	err = dao.InsertEmployer(&models.Employee{
		Identity:     id,
		Uid:          uid,
		Name:         e.Name,
		Birthday:     time.Unix(e.Birthday, 0),
		Phone:        e.Phone,
		Status:       0,
		Position:     e.Position,
		DepartmentId: e.Department,
	})
	//global.Global.Redis.Set(global.Global.Ctx, global.UidKey+use.Identity, use., 0)

	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.AddEmployerError)
		return
	}
	//添加进缓存
	go func() {
		_, err = global.Global.Redis.SAdd(global.Global.Ctx, global.Employer, id).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	result.Ok(c, nil)
}

// DeleteEmployee 删除员工
func DeleteEmployee(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//查询id是否存在
	e, err := dao.GetEmployerById(id)
	if err != nil || e.Uid == 0 {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//	删除员工信息
	err = dao.DeleteEmployer(id)
	if err != nil {
		result.Fail(c, global.BadRequest, global.DeleteError)
		return
	}
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
	e := new(employer)
	err := c.Bind(e)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryError)
		global.Global.Log.Error(err)
		return
	}
	//	更新信息

}

// EmployeeInfo 获取员工信息
func EmployeeInfo(c *gin.Context) {
	//获取identity
	identity := c.GetString("identity")
	if identity == "" {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if id == 0 || err != nil {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//获取员工信息
	val := global.Global.Redis.Get(global.Global.Ctx, global.Uid+identity).Val()
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
	info, err := dao.GetEmployerInfo(int64(id))
	if err != nil {
		result.Fail(c, global.DataNotFound, global.UserNotExistError)
		return
	}
	//插入
	go func() {
		if info == nil {
			_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+identity, "null", 0).Result()
			return
		}
		marshal, err := json.Marshal(info)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+identity, marshal, 0).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	result.Ok(c, info)
	return

}
