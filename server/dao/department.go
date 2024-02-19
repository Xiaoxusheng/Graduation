package dao

import (
	"server/global"
	"server/models"
)

func GetDepartmentList() ([]*models.Department, error) {
	list := make([]*models.Department, 0)
	err := global.Global.Mysql.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func DeleteDepartment(id string) error {
	d := new(models.Department)
	return global.Global.Mysql.Where("identity=?", id).Delete(d).Error
}
