package dao

import (
	"server/global"
	"server/models"
)

// GetMenuList 获取菜单列表
func GetMenuList() ([]models.Menu, error) {
	list := make([]models.Menu, 0)
	err := global.Global.Mysql.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// InsertMenu 增加菜单
func InsertMenu(menu *models.Menu) error {
	return global.Global.Mysql.Create(menu).Error
}

// DelMenu 删除菜单
func DelMenu(menu *models.Menu) error {
	return global.Global.Mysql.Unscoped().Where("menu_url=?", menu.MenuUrl).Delete(menu).Error
}

// UpdateMenu 更新 菜单信息
func UpdateMenu(menu *models.Menu) error {
	return global.Global.Mysql.Where("menu_url=?", menu.MenuUrl).Updates(menu).Error
}
