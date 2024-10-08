package dao

import (
	"fmt"
	"server/global"
	"server/models"
	"time"
)

// GetAttendanceList 根据uid获取考勤记录
func GetAttendanceList(limit, offset int, uid int64) ([]*models.Attendance, error) {
	list := make([]*models.Attendance, 0)
	err := global.Global.Mysql.Where("uid=?", uid).Limit(limit).Offset(offset - 1).Find(&list).Error
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
func GetDateList(limits, offset int, t time.Time) ([]global.ClockingIn, error) {
	list := make([]global.ClockingIn, 0)
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	//第二天0点
	t2 := t1.Add(time.Hour * 24)
	fmt.Println(t1, t2)
	err := global.Global.Mysql.Table("attendance_basic").Select("attendance_basic.uid,attendance_basic.Identity,attendance_basic.status,attendance_basic.start_time,attendance_basic.end_time,attendance_basic.date,employee_basic.department_id,employee_basic.name,employee_basic.sex").Joins("join employee_basic  on  employee_basic.uid=attendance_basic.uid").Where("attendance_basic.start_time>? and attendance_basic.end_time<?", t1, t2).Limit(limits).Offset(offset - 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateEndTime 加班
func UpdateEndTime(uid, time2 int64) error {
	attendances := new(models.Attendance)
	return global.Global.Mysql.Model(attendances).Where("uid=?", uid).Updates(&models.Attendance{
		EndTime: time.Unix(time2, 0),
		Status:  4,
	}).Error
}

/*
1表示缺勤 2表示打卡 3表示迟到 4表示加班 5表示补卡 6出差 7 请假
*/

// Leave 请假
func Leave(uid int64, id string, start, end time.Time) error {
	return global.Global.Mysql.Create(&models.Attendance{
		StartTime: start,
		EndTime:   end,
		Date:      start,
		Identity:  id,
		Uid:       uid,
		Name:      "",
		Status:    7,
	}).Error
}

// UpdateMakeCard 修改补卡
func UpdateMakeCard(uid int64, start, end time.Time, id string) error {
	return global.Global.Mysql.Create(&models.Attendance{
		Identity:  id,
		Uid:       uid,
		StartTime: start,
		EndTime:   end,
		Date:      start,
		Status:    5,
	}).Error
}

// ClockIn 员工上班打卡
func ClockIn(attendance *models.Attendance) error {
	return global.Global.Mysql.Create(attendance).Error
}

// AfterWork 下班打卡
func AfterWork(uid int64, t time.Time) error {
	attendance := new(models.Attendance)
	return global.Global.Mysql.Model(attendance).Where("uid=?  and  date=?", uid, t.Format(time.DateOnly)).Update("end_time", t).Error
}

func GetByStatus(date string) ([]models.Attendance, error) {
	list := make([]models.Attendance, 0)
	err := global.Global.Mysql.Where("date=?", date).Group("status").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
