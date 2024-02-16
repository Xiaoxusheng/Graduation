package dao

import (
	"server/global"
	"server/models"
)

func InsertEmployer(e *models.Employee) error {
	return global.Global.Mysql.Create(e).Error
}
