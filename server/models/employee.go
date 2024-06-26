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
	Sex          int32     `gorm:"type:int; comment:'1表是男员工, 2表示表示女 '" json:"sex,omitempty"`
	Status       int32     `gorm:"type:int; comment:'1表示正常在岗, 2表示离职 ,3表示开除'" json:"status,omitempty"`
	Position     int32     `gorm:"type:int; comment:职位 1 普通员工 2副主管 3主管 4 经理  5 经理 " json:"position,omitempty"`
	DepartmentId int32     `gorm:"type:int; comment:部门id" json:"department_id,omitempty"`
	IP           string    `gorm:"type:varchar(64) not null; comment:'IP地址'" json:"IP,omitempty"`
	ImageUrl     string    `gorm:"type:varchar(128); comment:'头像'" json:"image_url,omitempty"`
	gorm.Model
}

func (e *Employee) TableName() string {
	return "employee_basic"
}
