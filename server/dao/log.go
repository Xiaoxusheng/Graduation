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
func GetLogList(limits, offset int) ([]models.Log, error) {
	var count int64
	list := make([]models.Log, 0)
	err := global.Global.Mysql.Limit(limits).Offset(offset - 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	return list, nil
}

func GetCount(count *int64) error {
	return global.Global.Mysql.Table("log_basic").Count(count).Error
}
