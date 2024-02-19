package models

import "gorm.io/gorm"

// Examine 审核表
type Examine struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'打卡记录唯一标识'" json:"identity,omitempty"`
	Name     string `gorm:"type:varchar(10) not null unique; comment:'员工姓名'" json:"name,omitempty"`
	Uid      int64  `gorm:"type:varchar(36) not null unique; comment:'工号'" json:"uid,omitempty"`
	Status   int32  `gorm:"type int  not null ;comment:' 1表示加班申请 2表示补卡申请 3出差申请'" json:"status"`
	Reason   string `gorm:"type:varchar(100) not null unique; comment:'情况说明'" json:"reason,omitempty"`
	Pass     int32  `gorm:"type:int not null;comment:'是否通过 通过为1 不通过为0'" json:"pass"`
	gorm.Model
}

func (e *Examine) TableName() string {
	return "Examine_basic"
}
