package dao

import (
	"server/global"
	"server/models"
)

func InsertEmployer(e *models.Employee) error {
	return global.Global.Mysql.Create(e).Error
}

// GetEmployer 部门
func GetEmployer(department int32) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Where("department_id=?", department).Order("uid DESC").Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

// GetEmployerById 根据identity
func GetEmployerById(id string) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Where("identity=?", id).Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

// DeleteEmployer 删除
func DeleteEmployer(id string) error {
	e := new(models.Employee)
	return global.Global.Mysql.Where("identity=?", id).Delete(e).Error
}

// GetEmployerInfo uid查询
func GetEmployerInfo(uid int64) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Where("uid=?", uid).Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

func GetEmployerList(limit, offset int) ([]*models.Employee, error) {
	list := make([]*models.Employee, 0)
	err := global.Global.Mysql.Take("").Limit(limit).Offset(offset - 1).Find(&list).Error
	return list, err
}
