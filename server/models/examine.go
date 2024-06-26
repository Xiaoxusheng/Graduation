package models

import (
	"gorm.io/gorm"
	"time"
)

// Examine 审核表
type Examine struct {
	Identity  string    `gorm:"type:varchar(36) not null unique; comment:'打卡记录唯一标识'" json:"identity,omitempty"`
	Uid       int64     `gorm:"type:varchar(36) not null; comment:'工号'" json:"uid,omitempty"`
	StartTime time.Time `gorm:"type:datetime not null; comment:'开始时间'" json:"start_time,omitempty"`
	EndTime   time.Time `gorm:"type:datetime not null; comment:'结束时间'" json:"end_time,omitempty"`
	Status    int32     `gorm:"type int  not null ;comment:'1表示加班申请 2表示补卡申请 3出差申请 4请假申请'" json:"status"`
	Reason    string    `gorm:"type:varchar(100) not null ; comment:'情况说明'" json:"reason,omitempty"`
	Url       string    `gorm:"type varchar(16) not null;comment:'补卡照片'" json:"url"`
	Pass      int32     `gorm:"type:int not null;comment:'是否通过 通过为1 不通过为0'" json:"pass"`
	IsExamine int32     `gorm:"type:int not null;default:0; comment:'是否审核过 0表示未审核 1表示已经审核'" json:"is_examine,omitempty"`
	gorm.Model
}

func (e *Examine) TableName() string {
	return "examine_basic"
}
