package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
)

// AddDepartment 添加部门
func AddDepartment(c *gin.Context) {
	d := new(global.Department)
	err := c.Bind(d)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//	增加
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()
	err = dao.InsertDepartment(&models.Department{
		Identity: utils.Md5(d.Name),
		Name:     d.Name,
		Sort:     d.Sort,
		Leader:   d.Leader,
	})
	if err != nil {
		result.Fail(c, global.ServerError, global.AddDepartmentError)
		return
	}
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.DepartmentList)
	}()
	result.Ok(c, nil)

}

// DeleteDepartment 删除部门
func DeleteDepartment(c *gin.Context) {
	//获取删除的部门id
	id := c.Query("id")
	//删除部门信息,只有一人能操作
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()
	err := dao.DeleteDepartment(id)
	if err != nil {
		result.Fail(c, global.ResourceNotFound, global.DelDepartmentError)
		return
	}
	//	删除缓存
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.DepartmentList)
	}()
	result.Ok(c, nil)
}

// UpdateDepartment 更新部门信息
func UpdateDepartment(c *gin.Context) {
	d := new(global.Department)
	err := c.Bind(d)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//	拿到参数,验证
	global.Global.Mutex.Lock()
	err = dao.UpdateDepartment(d)
	global.Global.Mutex.Unlock()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.UpdateDepartmentError)
		return
	}
	//删除缓存
	go func() {
		global.Global.Redis.Del(global.Global.Ctx, global.DepartmentList)
	}()
	result.Ok(c, nil)
}

// GetDepartmentList 查询部门列表
func GetDepartmentList(c *gin.Context) {
	val := global.Global.Redis.Get(global.Global.Ctx, global.DepartmentList).Val()
	if val != "" {
		list := make([]models.Department, 0)
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			result.Fail(c, global.DataUnmarshal, global.DataUnmarshalError)
			return
		}
		result.Ok(c, list)
		return
	}
	global.Global.Mutex.RLock()
	list, err := dao.GetDepartmentList()
	global.Global.Mutex.RUnlock()
	if err != nil {
		result.Fail(c, global.ServerError, global.GetDepartmentError)
		return
	}
	result.Ok(c, list)
	//	同步
	go func() {
		marshal, err := json.Marshal(list)
		if err != nil {
			global.Global.Log.Warn(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.DepartmentList, marshal, global.DepartmentTime).Result()
		if err != nil {
			global.Global.Log.Warn(err)
		}
	}()
}
