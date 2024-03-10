package dao

import (
	"fmt"
	"server/global"
	"server/models"
)

// InsertLog 插入日志
func InsertLog(log *models.Log) error {
	return global.Global.Mysql.Create(log).Error
}

// GetLogList 查询日志
func GetLogList(s string, limits, offset int) ([]models.Log, error) {
	var count int64
	list := make([]models.Log, 0)
	err := global.Global.Mysql.Where("path like ?", s+"%").Limit(limits).Offset((offset - 1) * limits).Find(&list).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	return list, nil
}

func GetCount(s string, count *int64) error {
	return global.Global.Mysql.Table("log_basic").Where("path like ?", s+"%").Count(count).Error
}
