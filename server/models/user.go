package models

import "gorm.io/gorm"

// User 是用户表
type User struct {
	Username string `gorm:"type:varchar(10) not null unique; comment:'用户名'" json:"username,omitempty"`
	Identity string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty"`
	ImageUrl string `gorm:"type:varchar(64); comment:'头像'" json:"image_url,omitempty"`
	Password string `gorm:"type:varchar(64) not null; comment:'密码'" json:"password,omitempty"`
	Phone    string `gorm:"type:varchar(11) not null unique; comment:'手机号'" json:"phone,omitempty"`
	//Status   int32  `gorm:"type:int; comment:'0表示正常, 1表示封禁'" json:"status,omitempty"`
	IP      string `gorm:"type:varchar(64) not null; comment:'IP地址'" json:"IP,omitempty"`
	Salt    string `gorm:"type:varchar(30) not null; comment:'盐值'" json:"salt,omitempty"`
	Account int64  `gorm:"type:int not null unique; comment:'账号'" json:"account,omitempty"`
	gorm.Model
}

func (u *User) TableName() string {
	return "user_basic"
}
