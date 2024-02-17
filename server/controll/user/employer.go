package user

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

type employer struct {
	Username   string `json:"username,omitempty" `
	Identity   string `json:"identity,omitempty"`
	Uid        int32  `json:"uid,omitempty"`
	Name       string `json:"name,omitempty"`
	Birthday   int64  `json:"age,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Position   string `json:"position,omitempty"`
	Department int32  `json:"department,omitempty"`
}

// AddEmployee 添加员工信息
func AddEmployee(c *gin.Context) {
	e := new(employer)
	err := c.Bind(e)
	if err != nil {
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}

	//生成唯一员工标识
	uid, err := utils.GetUid(e.Department)
	if err != nil {
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}

	//	添加员工信息
	id := utils.GetUidV4()
	err = dao.InsertEmployer(&models.Employee{
		Username:   e.Username,
		Identity:   id,
		Uid:        uid,
		Name:       e.Name,
		Birthday:   time.Unix(e.Birthday, 0).Local(),
		Phone:      e.Phone,
		Status:     0,
		Position:   e.Position,
		Department: e.Department,
	})
	if err != nil {
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
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	limits, err := strconv.Atoi(limit)
	if err != nil {
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	list, err := dao.GetEmployerList(limits, offsets)
	if err != nil {
		result.Fail(c, global.BadRequest, global.GetEmployerListError)
		return
	}
	result.Ok(c, list)
}

// UpdateEmployee 更新员工信息
func UpdateEmployee(c *gin.Context) {

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
	val := global.Global.Redis.Get(global.Global.Ctx, identity).Val()
	if val != "" {
		e := new(employer)
		err = json.Unmarshal([]byte(val), e)
		if err != nil {
			return
		}
		result.Ok(c, e)
	}
	//不存在，过期
	info, err := dao.GetEmployerInfo(int64(id))
	if err != nil {
		result.Fail(c, global.DataNotFound, global.UserNotExistError)
		return
	}
	//插入
	go func() {
		if info.Name == "" {
			_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid, "null", 0).Result()
			return
		}
		marshal, err := json.Marshal(info)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid, marshal, 0).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	result.Ok(c, info)
	return

}
