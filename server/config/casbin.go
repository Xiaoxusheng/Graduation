package config

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"server/global"
)

func InitCasBin() {
	sql, err := gormadapter.NewAdapter("mysql", Config.Mysql.Username+":"+Config.Mysql.Password+"@tcp("+Config.Mysql.Url+")/"+Config.Mysql.Database+"?charset=utf8mb4&parseTime=True&loc=Local", true) // Your driver and data source.
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("./models.conf", sql)
	if err != nil {
		panic(err)
	}
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		panic(e)
	}

	global.Global.CasBin = e
	// Check the permission.

}
