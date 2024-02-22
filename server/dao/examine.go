package dao

import (
	"server/global"
	"server/models"
	"time"
)

// GetExamineList 请假审批表
func GetExamineList() ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? and created_at>?", 4, time.Now()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateLeaveStatus 请假审核
func UpdateLeaveStatus(uid, pass int32) error {
	examine := new(models.Examine)
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 4).Update("pass", pass).Error
}

// UpdateOvertimeStatus 加班申请
func UpdateOvertimeStatus(uid, pass int32) error {
	examine := new(models.Examine)
	//修改一下时间
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()

	err := UpdateEndTime(examine.EndTime)
	if err != nil {
		return err
	}
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 1).Updates(&models.Examine{
		Pass:      pass,
		IsExamine: 1,
	}).Error
}

// GetOvertimeList 获取加班申请列表,包括审批过的
func GetOvertimeList() ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? ", 1, time.Now()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
