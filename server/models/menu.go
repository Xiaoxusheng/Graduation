package models

import "gorm.io/gorm"

type Menu struct {
	MenuUrl       string `gorm:"type:varchar(100) not null unique; comment:'菜单路径'" json:"menuUrl,omitempty"`
	MenuName      string `gorm:"type:varchar(100) not null unique; comment:'菜单名称'" json:"menuName,omitempty"`
	Icon          string `gorm:"type:varchar(36); comment:'菜单图标'" json:"icon,omitempty"`
	ParentPath    string `gorm:"type:varchar(100); comment:'父级路径'" json:"parentPath,omitempty"`
	RouteName     string `gorm:"type:varchar(36) ; comment:'路由名称'" json:"routeName,omitempty"`
	Cacheable     bool   `gorm:"type:bool not null; default:true; comment:'是否缓存'" json:"cacheable,omitempty"`
	Badge         string `gorm:"type:varchar(36); comment:'角标'" json:"badge,omitempty"`
	LocalFilePath string `gorm:"type:varchar(36); comment:'本地路径'" json:"localFilePath,omitempty"`
	IsRootPath    bool   `gorm:"type:bool ; comment:'是否外链'" json:"isRootPath,omitempty"`
	Children      []Menu `gorm:"-"  json:"children"`
	gorm.Model    `json:"-"`
}

// Menu 表示菜单的数据结构
//type Menu struct {
//	MenuURL       string `gorm:"uniqueIndex;not null"` // 菜单URL，唯一且不能为空
//	MenuName      string `gorm:"not null"`             // 菜单名称，不能为空
//	RouteName     string // 路由名称
//	Icon          string // 菜单图标
//	ParentPath    string // 父级路径
//	IsRootPath    bool   // 是否是根路径
//	Badge         string // 角标
//	LocalFilePath string // 本地文件路径
//	Cacheable     bool   // 是否可缓存
//	IconPrefix    string // 图标前缀
//	IconFont      string // 图标字体
//	IconSelector  string // 图标选择器
//	IconPath      string // 图标路径
//	IsSingle      bool   // 是否是单例
//	Children      []Menu `gorm:"foreignKey:ParentPath;references:MenuURL"` // 子菜单，关联父级菜单的 MenuURL
//	gorm.Model
//}

func (m *Menu) TableName() string {
	return "menu_basic"
}
