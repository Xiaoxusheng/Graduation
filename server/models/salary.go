package models

import "gorm.io/gorm"

// Salary 工资表
type Salary struct {
	Identity                string  `gorm:"type:varchar(36) not null unique; comment:'工资记录唯一标识'" json:"identity,omitempty"`
	Uid                     int64   `gorm:"type:varchar(36) not null; comment:'工号'" json:"uid,omitempty"`
	PaidLeave               int32   `gorm:"type:int not null; comment:'带薪休假时长'" json:"paid_leave,omitempty"`
	Count                   int32   `gorm:"type:int not null; comment:'迟到早退次数,每次扣10元'" json:"count"`
	Total                   float64 `gorm:"type:float ; comment:'应发工资'" json:"total,omitempty"`
	ExpectedAttendanceHours float64 `gorm:"type:float not null; comment:'应出勤时长'" json:"expected_attendance_hours,omitempty"`
	AttendanceHours         float64 `gorm:"type:float not null; comment:'实际出勤时长'" json:"attendance_hours,omitempty"`
	Other                   float64 `gorm:"type:float not null; comment:'公积金，社保，五险一金等扣除的'" json:"other,omitempty"`
	Subsidy                 float64 `gorm:"type:float not null; comment:'公司补贴'" json:"subsidy,omitempty"`
	Date                    string  `gorm:"type:varchar(36)  not null; comment:'年-月'" json:"date,omitempty"`
	gorm.Model
}

func (e *Salary) TableName() string {
	return "salary_basic"
}
