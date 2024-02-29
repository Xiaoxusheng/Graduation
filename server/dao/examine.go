package dao

import (
	"server/global"
	"server/models"
	"time"
)

//'1表示加班申请 2表示补卡申请 3出差申请 4请假申请'

// GetExamineList 请假审批表
func GetExamineList(limits, offset int) ([]global.Applications, error) {
	list := make([]global.Applications, 0)
	err := global.Global.Mysql.Table("examine_basic").Select("employee_basic.sex,employee_basic.name,examine_basic.*").Joins("join employee_basic on   employee_basic.uid=examine_basic.uid").Where("examine_basic.status=? and examine_basic.created_at<?", 4, time.Now()).Limit(limits).Offset(offset - 1).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateLeaveStatus 请假审核
func UpdateLeaveStatus(uid int64, pass int32, id string, start, end time.Time) error {
	examine := new(models.Examine)
	//审核通过
	if pass == 1 {
		err := Leave(uid, id, start, end)
		if err != nil {
			return err
		}
	}
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 4).Update("pass", pass).Error
}

// UpdateOvertimeStatus 加班申请
func UpdateOvertimeStatus(uid int64, pass int32, endTime int64) error {
	examine := new(models.Examine)
	//修改一下时间
	global.Global.Mutex.Lock()
	defer global.Global.Mutex.Unlock()
	//审核通过
	if pass == 1 {
		err := UpdateEndTime(uid, endTime)
		if err != nil {
			return err
		}
	}
	return global.Global.Mysql.Model(examine).Where("uid=?  and status=?", uid, 1).Updates(&models.Examine{
		Pass:      pass,
		IsExamine: 1,
	}).Error
}

// GetOvertimeList 获取加班申请列表,包括审批过的
func GetOvertimeList(limits, offsets int) ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? ", 1).Limit(limits).Offset(offsets - 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// MakeCard 补卡申请审批
func MakeCard(uid int64, star, end time.Time, pass int32, id string) error {
	examine := new(models.Examine)
	if pass == 1 {
		err := UpdateMakeCard(uid, star, end, id)
		if err != nil {
			return err
		}
	}
	return global.Global.Mysql.Model(examine).Where("uid=?   and status=?", uid, 2).Updates(
		&models.Examine{
			Pass:      pass,
			IsExamine: 1,
		}).Error
}

// GetMarkCardList 补卡申请列表
func GetMarkCardList(limits, offsets int) ([]models.Examine, error) {
	list := make([]models.Examine, 0)
	err := global.Global.Mysql.Where("status=? ", 2).Limit(limits).Offset(offsets - 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// InsertMarkCardApplication 补卡申请
func InsertMarkCardApplication(examine *models.Examine) error {
	return global.Global.Mysql.Create(examine).Error
}

func GetByUid(uid int64, status int32) (*models.Examine, error) {
	examine := new(models.Examine)
	err := global.Global.Mysql.Where("uid=?  and status=? ", uid, status).First(examine).Error
	if err != nil {
		return nil, err
	}
	return examine, nil
}
