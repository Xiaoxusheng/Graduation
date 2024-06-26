package models

import "gorm.io/gorm"

type Notice struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'公告记录唯一标识'" json:"identity,omitempty"`
	Title    string `gorm:"type:varchar(36) not null; comment:'标题'" json:"title,omitempty"`
	Url      string `gorm:"type:varchar(200) not null; comment:'照片url'" json:"url,omitempty"`
	Text     string `gorm:"type:text not null; comment:'内容'" json:"text,omitempty"`
	Status   *int32 `gorm:"type:int not null; comment:'公告状态，1为显示，2为下架'" json:"status,omitempty"`
	Uid      int64  `gorm:"type:varchar(36) not null; comment:'创建者uid'" json:"uid,omitempty"`
	Date     int64  `gorm:"type:int not null; comment:'公告发布时间'" json:"date,omitempty"`
	gorm.Model
}

func (e *Notice) TableName() string {
	return "notice_basic"
}
