package dao

import (
	"server/global"
	"server/models"
)

func InsertLog(log *models.Log) error {
	return global.Global.Mysql.Create(log).Error
}
