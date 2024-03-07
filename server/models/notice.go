package models

import "gorm.io/gorm"

type Notice struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'公告记录唯一标识'" json:"identity,omitempty"`
	Uid      int64  `gorm:"type:varchar(36) not null; comment:'创建者uid'" json:"uid,omitempty"`
	Title    string `gorm:"type:varchar(36) not null; comment:'标题'" json:"title,omitempty"`
	Text     string `gorm:"type:varchar(36) not null; comment:'内容'" json:"text,omitempty"`
	Status   *int32 `gorm:"type:int not null; comment:'公告状态，1为显示，2为下架'" json:"status,omitempty"`
	Date     int64  `gorm:"type:int not null; comment:'公告发布时间'" json:"date,omitempty"`
	gorm.Model
}

func (e *Notice) TableName() string {
	return "notice_basic"
}
