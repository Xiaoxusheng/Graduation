package dao

import (
	"server/global"
	"server/models"
)

func GetMenuList() ([]models.Menu, error) {
	list := make([]models.Menu, 0)
	err := global.Global.Mysql.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
