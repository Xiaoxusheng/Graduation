package models

import "gorm.io/gorm"

// Log 日志表
type Log struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty"`
	Uid      string `gorm:"type:varchar(36) not null unique; comment:'请求人的uid'" json:"uid,omitempty"`
	Method   string `gorm:"type:varchar(10) not null; comment:'请求方式'" json:"method,omitempty"`
	HttpCode int32  `gorm:"type:int not null; comment:'http状态码'" json:"http_code,omitempty"`
	Path     string `gorm:"type:varchar(36) not null; comment:'资源路径'" json:"path,omitempty"`
	gorm.Model
}

func (l *Log) TableName() string {
	return "log_basic"
}
