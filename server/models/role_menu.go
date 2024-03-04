package models

import "gorm.io/gorm"

type RoleMenu struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'工资记录唯一标识'" json:"identity,omitempty"`
	Role     string `gorm:"type:varchar(36) not null unique; comment:'角色'" json:"role,omitempty"`
	Menu     string `gorm:"type:varchar(400) not null unique; comment:'菜单的id列表'" json:"menu,omitempty"`
	gorm.Model
}

func (e *RoleMenu) TableName() string {
	return "role_menu_basic"
}
