package dao

import (
	"server/global"
	"server/models"
)

// GetSalary 获取单人
func GetSalary(uid int64, date string) (*models.Salary, error) {
	salary := new(models.Salary)
	err := global.Global.Mysql.Where("uid=? and date=?", uid, date).First(salary).Error
	if err != nil {
		return nil, err
	}
	return salary, nil
}

// InsertSalary 添加
func InsertSalary(salary *models.Salary) error {
	return global.Global.Mysql.Create(salary).Error
}

// UpdateSalary 更新
func UpdateSalary(uid int64, count int32, h float64) error {
	return global.Global.Mysql.Model(new(models.Salary)).Where("uid=?", uid).Updates(&models.Salary{
		Count:           count,
		AttendanceHours: h,
	}).Error
}

func GetSalaryList(limits, offset int, t string) ([]*models.Salary, error) {
	list := make([]*models.Salary, 0)
	err := global.Global.Mysql.Where("date=?", t).Limit(limits).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
