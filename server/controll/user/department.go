package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
)

// AddDepartment 添加部门
func AddDepartment(c *gin.Context) {

}

// DeleteDepartment 删除部门
func DeleteDepartment(c *gin.Context) {
	//获取删除的部门id
	id := c.Query("id")
	//删除部门信息
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
	d := new(models.Department)
	err := c.Bind(d)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//	拿到参数

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
	list, err := dao.GetDepartmentList()
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
