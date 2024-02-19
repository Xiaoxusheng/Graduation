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

func UpdateDepartment(department *global.Department) error {
	var l, name, sort string
	if department.Leader != "" {
		l = "leader"
	}
	if department.Name != "" {
		name = "name"
	}
	if department.Sort != 0 {
		sort = "sort"
	}
	return global.Global.Mysql.Model(department).Select(l, name, sort).Updates(department).Error
}

func InsertDepartment(department *models.Department) error {
	return global.Global.Mysql.Create(department).Error
}
