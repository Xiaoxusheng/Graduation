package admin

import (
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/result"
)

// LeaveApplication 请假申请审核
func LeaveApplication(c *gin.Context) {
	application := new(global.Application)
	err := c.Bind(application)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	//判断员工uid是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, application.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	if application.Pass == 0 || application.Pass == 1 {
		err = dao.UpdateLeaveStatus(int32(application.Uid), application.Pass)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.LeaveApplicationError)
			return
		}
		result.Ok(c, nil)
		return
	}
	result.Fail(c, global.BadRequest, global.QueryError)
}

// GetLeaveApplicationList 获取请假申请列表
func GetLeaveApplicationList(c *gin.Context) {
	list, err := dao.GetExamineList()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetLeaveListError)
		return
	}
	result.Ok(c, list)
}

// OvertimeApplication 加班申请审核
func OvertimeApplication(c *gin.Context) {
	application := new(global.Application)
	err := c.Bind(application)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	//判断员工uid是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, application.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	if application.Pass == 0 || application.Pass == 1 {
		err = dao.UpdateOvertimeStatus(int32(application.Uid), application.Pass)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.OverTimeApplicationError)
			return
		}
		result.Ok(c, nil)
		return
	}

}

// GetOvertimeList 获取加班申请审核列表
func GetOvertimeList(c *gin.Context) {
	list, err := dao.GetOvertimeList()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetOverTimeError)
		return
	}
	result.Ok(c, list)

}

// MakeCardApplication 补卡申请审核
func MakeCardApplication(c *gin.Context) {
	application := new(global.Application)
	err := c.Bind(application)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	//判断员工uid是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, application.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//
	if application.Pass == 0 || application.Pass == 1 {
		err = dao.MakeCard(int32(application.Uid), application.Pass)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.MarkCardApplicationError)
			return
		}
		result.Ok(c, nil)
		return
	}
}

// GetMarkCardList 补卡申请列表
func GetMarkCardList(c *gin.Context) {
	list, err := dao.GetMarkCardList()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetMarkCardLiatError)
		return
	}
	result.Ok(c, list)
}
