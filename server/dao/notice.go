package dao

import (
	"server/global"
	"server/models"
)

func InsertNotice(notice *models.Notice) error {
	return global.Global.Mysql.Create(notice).Error
}

func UpdateNotice(id string, status int32) error {
	return global.Global.Mysql.Model(new(models.Notice)).Where("identity=?", id).Update("status", status).Error
}

// 用户

func GetNoticeList() ([]*models.Notice, error) {
	list := make([]*models.Notice, 0)
	err := global.Global.Mysql.Where("status=?", 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
