package dao

import (
	"server/global"
	"server/models"
)

func InsertRoleMenu(menu *models.RoleMenu) error {
	return global.Global.Mysql.Create(menu).Error
}

func GetMenuLists(role string) ([]models.Menu, error) {
	list := make([]models.Menu, 0)
	err := global.Global.Mysql.Table("menu_basic").Where("uid in (?)", global.Global.Mysql.Model(new(models.RoleMenu)).Select("menu").Where("role=?", role)).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}