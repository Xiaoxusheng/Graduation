package dao

import (
	"server/global"
	"server/models"
	"time"
)

// GetAttendanceList 根据uid获取考勤记录
func GetAttendanceList(limit, offset int, uid int32) ([]*models.Attendance, error) {
	list := make([]*models.Attendance, 0)
	err := global.Global.Mysql.Where("uid=?", uid).Limit(limit).Offset(offset).Find(&list).Error
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
	}).Error
}

// GetDateList 获取某一天所有的打卡记录
func GetDateList(limits, offset int, t time.Time) ([]models.Attendance, error) {
	list := make([]models.Attendance, 0)
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	//第二天0点
	t2 := t1.Add(time.Hour * 24)
	err := global.Global.Mysql.Where("start_time>? and end_time<?", t1, t2).Limit(limits).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateEndTime 加班
func UpdateEndTime(uid int32, time2 int64) error {
	attendances := new(models.Attendance)
	return global.Global.Mysql.Model(attendances).Where("uid=?", uid).Updates(&models.Attendance{
		EndTime: time.Unix(time2, 0),
		Status:  4,
	}).Error
}

// UpdateMakeCard 修改补卡
func UpdateMakeCard(uid int32) error {
	attendances := new(models.Attendance)
	return global.Global.Mysql.Model(attendances).Where("uid=?", uid).Updates(&models.Attendance{
		Status: 5,
	}).Error
}

// ClockIn 员工上班打卡
func ClockIn(attendance *models.Attendance) error {
	return global.Global.Mysql.Create(attendance).Error
}

// AfterWork 下班打卡
func AfterWork(uid int64, t int64) error {
	attendance := new(models.Attendance)
	return global.Global.Mysql.Model(attendance).Where("uid=?", uid).Update("end_time", t).Error
}
