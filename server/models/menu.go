package models

import "gorm.io/gorm"

// Menu 菜单表
type Menu struct {
	Uid           int32  `gorm:"type:varchar(100) not null unique; comment:'菜单路径'" json:"uid"`
	MenuUrl       string `gorm:"type:varchar(100) not null unique; comment:'菜单路径'" json:"menuUrl,omitempty"`
	MenuName      string `gorm:"type:varchar(100) not null unique; comment:'菜单名称'" json:"menuName,omitempty"`
	Icon          string `gorm:"type:varchar(36); comment:'菜单图标'" json:"icon,omitempty"`
	ParentPath    string `gorm:"type:varchar(100); comment:'父级路径'" json:"parentPath,omitempty"`
	RouteName     string `gorm:"type:varchar(36) ; comment:'路由名称'" json:"routeName,omitempty"`
	Cacheable     bool   `gorm:"type:bool not null; default:true; comment:'是否缓存'" json:"cacheable,omitempty"`
	Badge         string `gorm:"type:varchar(36); comment:'角标'" json:"badge,omitempty"`
	LocalFilePath string `gorm:"type:varchar(36); comment:'本地路径'" json:"localFilePath,omitempty"`
	IsRootPath    bool   `gorm:"type:bool ; comment:'是否外链'" json:"isRootPath,omitempty"`
	Hidden        bool   `gorm:"type:bool ; comment:'是否隐藏'" json:"hidden,omitempty"`
	Children      []Menu `gorm:"-"  json:"children"`
	gorm.Model    `json:"-"`
}

func (m *Menu) TableName() string {
	return "menu_basic"
}
