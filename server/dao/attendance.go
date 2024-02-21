package dao

import (
	"server/global"
	"server/models"
	"time"
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

// UpdateAttendance 更新考勤数据
func UpdateAttendance(attendance *global.Attendance) error {
	return global.Global.Mysql.Where("uid=?", attendance.Uid).Updates(&models.Attendance{
		Identity:  attendance.Identity,
		Uid:       attendance.Uid,
		Name:      attendance.Name,
		StartTime: time.Unix(attendance.StartTime, 0),
		EndTime:   time.Unix(attendance.EndTime, 0),
		Status:    attendance.Status,
		Reason:    attendance.Reason,
	}).Error
}
