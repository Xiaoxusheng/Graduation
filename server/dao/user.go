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

func GetUserAccountPwd(account, password string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("account=? and password=?", account, password).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetInfoByIdentity(id string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("identity=?", id).Select("username", "identity", "account", "image_url", "password", "phone", "IP", "ID", "CreatedAt", "UpdatedAt").Take(user).Error
	return user, err
}

func GetUsername(username string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("username=?", username).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetPhone(phone string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("phone=?", phone).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// InsertUser 新增用户
func InsertUser(user *models.User) error {
	return global.Global.Mysql.Create(user).Error
}

// UpdatePwd 重置密码
func UpdatePwd(account string, pwd string) error {
	user := new(models.User)
	return global.Global.Mysql.Model(user).Where("account=?", account).Update("password", pwd).Error
}

func GetByAccount(account string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("account=?", account).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除
func DeleteUser(id string) error {
	user := new(models.User)
	return global.Global.Mysql.Where("identity=?", id).Delete(user).Error
}

// 修改密码
func UpdatePwds() {

}
