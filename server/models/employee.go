package models

import (
	"gorm.io/gorm"
	"time"
)

// Employee 员工信息表
type Employee struct {
	Username   string    `gorm:"type:varchar(10);not null unique; comment:'用户名'" json:"username,omitempty"`
	Identity   string    `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty"`
	Uid        int32     `gorm:"type:varchar(36) not null unique; comment:'工号'" json:"uid,omitempty"`
	Name       string    `gorm:"type:varchar(10) not null unique; comment:'员工姓名'" json:"name,omitempty"`
	Age        time.Time `gorm:"type:time not null unique; comment:'员工生日'" json:"age,omitempty"`
	Phone      string    `gorm:"type:varchar(11) not null unique; comment:'手机号'" json:"phone,omitempty"`
	Status     int32     `gorm:"type:int; comment:'0表示正常在岗, 1表示离职 ,表示开除'" json:"status,omitempty"`
	Position   string    `gorm:"type:int; comment:职位" json:"position,omitempty"`
	Department int32     `gorm:"type:int unique; comment:部门" json:"department,omitempty"`
	gorm.Model
}

func (e *Employee) TableName() string {
	return "employee_basic"
}
