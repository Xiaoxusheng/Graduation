package user

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
)

// AddEmployee 添加员工信息
func AddEmployee(c *gin.Context) {

}

// DeleteEmployee 删除员工
func DeleteEmployee(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		result.Fail(c, global.BadRequest, global.QueryNotFound)
		return
	}
	//查询id是否存在

}

// EmployeeList 员工列表
func EmployeeList(c *gin.Context) {

}

// UpdateEmployee 更新员工信息
func UpdateEmployee(c *gin.Context) {

}

// EmployeeInfo 获取员工信息
func EmployeeInfo(c *gin.Context) {

}
