package dao

import (
	"server/global"
	"server/models"
)

// GetDepartmentList  获取部门列表
func GetDepartmentList() ([]*models.Department, error) {
	list := make([]*models.Department, 0)
	err := global.Global.Mysql.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteDepartment 删除部门
func DeleteDepartment(id string) error {
	d := new(models.Department)
	//return global.Global.Mysql.Model(d).Where("identity=?", id).Update("status", 1).Error
	return global.Global.Mysql.Unscoped().Where("identity=?", id).Delete(d).Error
}

// StopDepartment 停用部门
func StopDepartment(id string) error {
	d := new(models.Department)
	return global.Global.Mysql.Model(d).Where("identity=?", id).Update("status", 1).Error
}

// UpdateDepartment 更新部门信息
func UpdateDepartment(department *global.Department) error {
	d := new(models.Department)
	return global.Global.Mysql.Model(d).Where("identity=?", department.Identity).Updates(&models.Department{
		Status: department.Status,
		Name:   department.Name,
		Sort:   department.Sort,
		Leader: department.Leader,
	}).Error
}

// InsertDepartment 增加部门
func InsertDepartment(department *models.Department) error {
	return global.Global.Mysql.Create(department).Error
}
