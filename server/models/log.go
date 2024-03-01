package models

import (
	"gorm.io/gorm"
)

// Log 日志表
type Log struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty"`
	Method   string `gorm:"type:varchar(10) not null; comment:'请求方式'" json:"method,omitempty"`
	Path     string `gorm:"type:varchar(36) not null; comment:'资源路径'" json:"path,omitempty"`
	IP       string `gorm:"type:varchar(36) not null; comment:'ip地址'" json:"ip,omitempty"`
	Time     int64  `gorm:"type:float not null ; comment:'耗时'" json:"time,omitempty"`
	Uid      int64  `gorm:"type:varchar(36) not null; comment:'请求人的uid'" json:"uid,omitempty"`
	HttpCode int32  `gorm:"type:int not null; comment:'http状态码'" json:"http_code,omitempty"`
	gorm.Model
}

func (l *Log) TableName() string {
	return "log_basic"
}
