package dao

import (
	"server/global"
	"server/models"
	"time"
)

//1表示缺勤 2表示打卡 3表示迟到 4表示加班 5表示补卡 6出差 7 请假

// GetExamineList 请假审批表
func GetExamineList() ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? and created_at>?", 7, time.Now()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateLeaveStatus 请假审核
func UpdateLeaveStatus(uid, pass int32) error {
	examine := new(models.Examine)
	if pass == 1 {

	}
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 7).Update("pass", pass).Error
}

// UpdateOvertimeStatus 加班申请
func UpdateOvertimeStatus(uid, pass int32) error {
	examine := new(models.Examine)
	//修改一下时间
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()

	err := UpdateEndTime(uid, examine.EndTime)
	if err != nil {
		return err
	}
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 4).Updates(&models.Examine{
		Pass:      pass,
		IsExamine: 1,
	}).Error
}

// GetOvertimeList 获取加班申请列表,包括审批过的
func GetOvertimeList() ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? ", 4).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// MakeCard 补卡申请审批
func MakeCard(uid, pass int32) error {
	examine := new(models.Examine)
	if pass == 1 {
		err := UpdateMakeCard(uid)
		if err != nil {
			return err
		}
	}
	return global.Global.Mysql.Model(examine).Where("uid=?   and status=?", uid, 5).Updates(
		&models.Examine{
			Pass:      pass,
			IsExamine: 1,
		}).Error
}

// GetMarkCardList 补卡申请列表
func GetMarkCardList() ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? ", 5).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// InsertMarkCardApplication 补卡申请
func InsertMarkCardApplication(examine *models.Examine) error {
	return global.Global.Mysql.Create(examine).Error
}
