package dao

import (
	"server/global"
	"server/models"
)

func InsertNotice(notice *models.Notice) error {
	return global.Global.Mysql.Create(notice).Error
}

func UpdateNotice(notice *global.UpdateNotice) error {
	return global.Global.Mysql.Model(new(models.Notice)).Where("identity=?", notice.Id).Updates(map[string]any{
		"status": notice.Status,
		"title":  notice.Title,
		"text":   notice.Text,
		"url":    notice.Url,
	}).Error
}

func GetNoticeLists() ([]*models.Notice, error) {
	list := make([]*models.Notice, 0)
	err := global.Global.Mysql.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetNoticeList 用户只拿到没下架的
func GetNoticeList() ([]*models.Notice, error) {
	list := make([]*models.Notice, 0)
	err := global.Global.Mysql.Where("status=?", 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func DeleteNotice(id string) error {
	return global.Global.Mysql.Where("identity=?", id).Delete(new(models.Notice)).Error
}

func GetExists(id string) bool {
	notice := new(models.Notice)
	err := global.Global.Mysql.Where("identity=?", id).Take(notice).Error
	if err != nil || notice.Text == "" {
		global.Global.Log.Error(err)
		return false
	}
	return true
}
