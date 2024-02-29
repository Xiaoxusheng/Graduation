package admin

import (
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/result"
	"server/utils"
	"strconv"
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
	if application.Pass == 1 || application.Pass == 2 {
		info, err := dao.GetByUid(application.Uid, 1)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.OverTimeApplicationError)
			return
		}

		id := utils.GetUidV4()
		err = dao.UpdateLeaveStatus(application.Uid, application.Pass, id, info.StartTime, info.EndTime)
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
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")

	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	list, err := dao.GetExamineList(limits, offsets)
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
	if application.Pass == 1 || application.Pass == 2 {
		info, err := dao.GetByUid(application.Uid, 1)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.OverTimeApplicationError)
			return
		}
		err = dao.UpdateOvertimeStatus(application.Uid, application.Pass, info.EndTime.Unix())
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
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")

	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	list, err := dao.GetOvertimeList(limits, offsets)
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
	if application.Pass == 1 || application.Pass == 2 {
		info, err := dao.GetByUid(application.Uid, 2)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.MarkCardApplicationError)
			return
		}
		id := utils.GetUidV4()
		err = dao.MakeCard(application.Uid, info.StartTime, info.EndTime, application.Pass, id)
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
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "1")

	limits, err := strconv.Atoi(limit)
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AtoiError)
		return
	}
	list, err := dao.GetMarkCardList(limits, offsets)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetMarkCardLiatError)
		return
	}
	result.Ok(c, list)
}
