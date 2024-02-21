package dao

import (
	"server/global"
	"server/models"
	"time"
)

func InsertEmployer(e *models.Employee, fn func()) error {
	err := global.Global.Mysql.Create(e).Error
	if err != nil {
		fn()
		return err
	}
	return nil
}

// GetEmployer  根据部门id获取员工
func GetEmployer(department int32) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Where("department_id=?", department).Order("uid DESC").Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

func GetEmployerByUid(uid string) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Where("uid=?", uid).Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

// DeleteEmployer 删除
func DeleteEmployer(uid string) error {
	e := new(models.Employee)
	return global.Global.Mysql.Unscoped().Where("uid=?", uid).Delete(e).Error
}

// GetEmployerInfo uid查询
func GetEmployerInfo(uid int64) (*models.Employee, error) {
	e := new(models.Employee)
	err := global.Global.Mysql.Select("user_basic.*", "employee_basic.*").Table("employee_basic").Joins("join user_basic on user_basic.identity=employee_basic.identity").Where("employee_basic.uid=?", uid).Take(&e).Error
	//err := global.Global.Mysql.Where("uid=?", uid).Take(e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

// GetEmployerList 获取员工列表
func GetEmployerList(limit, offset int) ([]*global.UserInfo, error) {
	list := make([]*global.UserInfo, 0)
	//err := global.Global.Mysql.Select("user_basic.*", "employee_basic.*").Table("employee_basic").Joins("join user_basic on user_basic.identity=employee_basic.identity").Limit(limit).Offset(offset - 1).Find(&list).Error
	err := global.Global.Mysql.Table("employee_basic").Limit(limit).Offset(offset - 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

// UpdateEmployer 更新员工信息
func UpdateEmployer(employers *global.Employers) error {
	return global.Global.Mysql.Where("uid=?", employers.Uid).Updates(&models.Employee{
		Uid:          employers.Uid,
		Name:         employers.Name,
		Birthday:     time.Unix(employers.Birthday, 0),
		Phone:        employers.Phone,
		Status:       employers.Status,
		Position:     employers.Position,
		DepartmentId: employers.Department,
	}).Error
}
