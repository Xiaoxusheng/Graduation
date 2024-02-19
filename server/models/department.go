package models

import "gorm.io/gorm"

// Department 部门表
type Department struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty"`
	Name     string `gorm:"type:varchar(10);not null unique; comment:'部门名称'" json:"name,omitempty"`
	Sort     int32  `gorm:"type:int not null unique; comment:'部门编号'" json:"sort,omitempty"`
	Leader   string `gorm:"type:varchar(36);not null unique; comment:'主管'" json:"leader,omitempty"`
	gorm.Model
}

func (d *Department) TableName() string {
	return "department_basic"
}
