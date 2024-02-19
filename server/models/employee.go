package models

import (
	"gorm.io/gorm"
	"time"
)

// Employee 员工信息表
type Employee struct {
	Identity     string    `gorm:"type:varchar(36) not null unique; comment:'员工唯一标识'" json:"identity,omitempty"`
	Uid          int64     `gorm:"type:varchar(36) not null unique; comment:'工号'" json:"uid,omitempty"`
	Name         string    `gorm:"type:varchar(10) not null unique; comment:'员工姓名'" json:"name,omitempty"`
	Birthday     time.Time `gorm:"type:datetime not null unique; comment:'员工生日'" json:"birthday,omitempty"`
	Phone        string    `gorm:"type:varchar(11) not null unique; comment:'手机号'" json:"phone,omitempty"`
	Status       int32     `gorm:"type:int; comment:'0表示正常在岗, 1表示离职 ,表示开除'" json:"status,omitempty"`
	Position     int32     `gorm:"type:int; comment:职位 0 普通员工 1副主管 2主管 3 经理  4 经理 " json:"position,omitempty"`
	DepartmentId int32     `gorm:"type:int unique; comment:部门id" json:"department_id,omitempty"`
	gorm.Model
}

func (e *Employee) TableName() string {
	return "employee_basic"
}
