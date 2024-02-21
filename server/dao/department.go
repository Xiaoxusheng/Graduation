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
	return global.Global.Mysql.Unscoped().Where("identity=?", id).Delete(d).Error
}

func UpdateDepartment(department *global.Department) error {
	d := new(models.Department)
	return global.Global.Mysql.Model(d).Where("identity=?", department.Identity).Updates(&models.Department{
		Name:   department.Name,
		Sort:   department.Sort,
		Leader: department.Leader,
	}).Error
}

func InsertDepartment(department *models.Department) error {
	return global.Global.Mysql.Create(department).Error
}
