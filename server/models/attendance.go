package models

import (
	"gorm.io/gorm"
	"time"
)

// Attendance 考勤表
type Attendance struct {
	Identity     string    `gorm:"type:varchar(36) not null unique; comment:'打卡记录唯一标识'" json:"identity,omitempty"`
	Uid          int64     `gorm:"type:varchar(36) not null unique; comment:'工号'" json:"uid,omitempty"`
	Name         string    `gorm:"type:varchar(10) not null; comment:'员工姓名'" json:"name,omitempty"`
	DepartmentId int32     `gorm:"type:int not null; comment:'员工部门'" json:"department_id,omitempty"`
	StartTime    time.Time `gorm:"type:datetime not null; comment:'上班前打卡时间'" json:"start_time,omitempty"`
	EndTime      time.Time `gorm:"type:datetime not null; comment:'下班后打卡时间'" json:"end_time,omitempty"`
	Status       int32     `gorm:"type int  not null default:1;comment:'1表示缺勤 2表示打卡 3表示迟到 4表示加班 5表示补卡 6出差 7 请假'" json:"status"`
	gorm.Model
}

func (l *Attendance) TableName() string {
	return "attendance_basic"
}
