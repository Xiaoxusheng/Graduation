package models

import (
	"gorm.io/gorm"
	"time"
)

// Attendance 考勤表
type Attendance struct {
	Identity string        `gorm:"type:varchar(36) not null unique; comment:'打卡记录唯一标识'" json:"identity,omitempty"`
	Username string        `gorm:"type:varchar(10);not null unique; comment:'用户名'" json:"username,omitempty"`
	Uid      int64         `gorm:"type:varchar(36) not null unique; comment:'工号'" json:"uid,omitempty"`
	Name     string        `gorm:"type:varchar(10) not null unique; comment:'员工姓名'" json:"name,omitempty"`
	Duration time.Duration `gorm:"type int not null;comment:'员工上班时长'"  json:"duration"`
	Status   int32         `gorm:"type int  not null ;comment:'0表示缺勤 1表示打卡 2表示迟到  3表示加班 5表示补卡'"`
	Url      string        `gorm:"type varchar(16) not null;comment:'补卡照片'"`
	gorm.Model
}

func (l *Attendance) TableName() string {
	return "attendance_basic"
}
