package dao

import (
	"server/global"
	"server/models"
)

// GetAttendanceList 根据uid获取考勤记录
func GetAttendanceList(uid int32) ([]*models.Attendance, error) {
	list := make([]*models.Attendance, 0)
	err := global.Global.Mysql.Where("uid=?", uid).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
