package dao

import (
	"errors"
	"gorm.io/gorm"
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

func GetInfoByIdentity(id string) (*global.AdminInfo, error) {
	user := new(global.AdminInfo)
	err := global.Global.Mysql.Table("user_basic").Select("user_basic.identity,user_basic.username,employee_basic.*").Joins("join employee_basic on user_basic.account=employee_basic.uid").Where("user_basic.identity=?", id).Scan(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
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

func GetUid(identity string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("identity=?", identity).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(tx *gorm.DB, id, name string) error {
	return tx.Model(new(models.User)).Where("identity=?", id).Update("username", name).Error
}

func GetSalt(username, password string) (*models.User, error) {
	user := new(models.User)
	err := global.Global.Mysql.Where("username=? and password=?", username, password).First(user).Error
	if err != nil {
		return nil, err
	}
	if user.Salt == "" {
		return nil, errors.New("用户不存在！")
	}
	return user, nil
}
