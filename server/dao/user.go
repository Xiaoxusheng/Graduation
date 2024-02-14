package dao

import (
	"server/global"
	"server/models"
)

func GetUserByUsePwd(username, password string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("username=? and password=?", username, password).Take(user).Error
	return user, err
}

func GetInfoByIdentity(id string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("identity=?", id).Take(user).Error
	return user, err
}
