package config

import (
	"log"
	"os"
	"server/global"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//数据库

var (
	onc sync.Once
	db  *gorm.DB
	err error
)

func InitMysql() {
	onc.Do(
		func() {
			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Info, // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      true,        // Don't include params in the SQL log
					Colorful:                  false,       // Disable color
				},
			)
			dsn := Config.Mysql.Username + ":" + Config.Mysql.Password + "@tcp(" + Config.Mysql.Url + ")/" + Config.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger:      newLogger,
				PrepareStmt: true,
			})
			if err != nil {
				panic(err)
			}
			db.Debug()
			global.Global.Log.Info("mysql连接成功！")
			mysqlDB, err := db.DB()
			if err != nil {
				panic(err)
			}
			// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
			mysqlDB.SetMaxIdleConns(Config.Mysql.MaxIdleCons)

			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			mysqlDB.SetMaxOpenConns(Config.Mysql.MaxOpenCons)

			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			mysqlDB.SetConnMaxLifetime(time.Minute * time.Duration(Config.Mysql.ConnMaxLifetime))
			global.Global.Mysql = db
		})
	//建表
	//err = global.Global.Mysql.AutoMigrate(&models.Notice{})
	//err = global.Global.Mysql.AutoMigrate(&models.RoleMenu{})
	//err = global.Global.Mysql.AutoMigrate(&models.User{})
	//err = global.Global.Mysql.AutoMigrate(&models.Department{})
	//err = global.Global.Mysql.AutoMigrate(&models.Log{})
	//err = global.Global.Mysql.AutoMigrate(&models.Attendance{})
	//err = global.Global.Mysql.AutoMigrate(&models.Menu{})
	//err = global.Global.Mysql.AutoMigrate(&models.Examine{})
	//err = global.Global.Mysql.AutoMigrate(&models.Employee{})
	//err = global.Global.Mysql.AutoMigrate(&models.Salary{})
	//if err != nil {
	//	global.Global.Log.Error(err)
	//}

	//err = global.Global.Mysql.Create(&models.Menu{
	//	MenuUrl:       "/other",
	//	MenuName:      "功能/组件",
	//	Icon:          "icon-apps",
	//	ParentPath:    "",
	//	RouteName:     "",
	//	Badge:         "",
	//	LocalFilePath: "",
	//}).Error
	//global.Global.Log.Warn(err)
	//err = global.Global.Mysql.Create([]models.Menu{
	//	{
	//		MenuUrl:    "/index",
	//		MenuName:   "Dashborad",
	//		Icon:       "icon-dashboard",
	//		ParentPath: "",
	//		RouteName:  "dashborad",
	//		Cacheable:  true,
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/index/home",
	//		MenuName:   "主控台",
	//		RouteName:  "home",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/index",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/index/work-place",
	//		MenuName:   "工作台",
	//		RouteName:  "workPlace",
	//		Icon:       "icon-subscribed",
	//		Cacheable:  true,
	//		ParentPath: "/index",
	//		Badge:      "",
	//		IsRootPath: true,
	//	},
	//	{
	//		MenuUrl:    "/system",
	//		MenuName:   "系统管理",
	//		RouteName:  "system",
	//		Icon:       "icon-settings",
	//		Cacheable:  true,
	//		ParentPath: "",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:       "/systems/department",
	//		MenuName:      "部门管理",
	//		RouteName:     "department",
	//		Icon:          "icon-apps",
	//		Cacheable:     true,
	//		ParentPath:    "/systems",
	//		Badge:         "new",
	//		LocalFilePath: "/system/local-path/department",
	//	},
	//	{
	//		MenuUrl:    "/systems/user",
	//		MenuName:   "员工管理",
	//		RouteName:  "user",
	//		Icon:       "icon-user",
	//		Cacheable:  true,
	//		ParentPath: "/systems",
	//		Badge:      "dot",
	//	},
	//	{
	//		MenuUrl:    "/systems/role",
	//		MenuName:   "角色管理",
	//		RouteName:  "",
	//		Icon:       "icon-stamp",
	//		Cacheable:  true,
	//		ParentPath: "/systems",
	//		Badge:      "12",
	//	},
	//	{
	//		MenuUrl:    "/systems/menu",
	//		MenuName:   "菜单管理",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/systems",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/list",
	//		MenuName:   "考勤管理",
	//		RouteName:  "",
	//		Icon:       "icon-list",
	//		Cacheable:  true,
	//		ParentPath: "",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/list/table-with-search",
	//		MenuName:   "请假申请",
	//		RouteName:  "",
	//		Icon:       "icon-user-group",
	//		Cacheable:  true,
	//		ParentPath: "/list",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/list/table-custom",
	//		MenuName:   "加班申请",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/list",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/list/list",
	//		MenuName:   "补卡申请",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/list",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/list/card-list",
	//		MenuName:   "考勤列表",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/list",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/form",
	//		MenuName:   "表单demo",
	//		RouteName:  "",
	//		Icon:       "icon-edit",
	//		Cacheable:  true,
	//		ParentPath: "",
	//		Badge:      "dot",
	//	},
	//	{
	//		MenuUrl:    "/form/base-form-view",
	//		MenuName:   "表单",
	//		RouteName:  "",
	//		Icon:       "icon-apps",
	//		Cacheable:  true,
	//		ParentPath: "/form",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:       "/other",
	//		MenuName:      "功能/组件",
	//		Icon:          "icon-apps",
	//		ParentPath:    "",
	//		RouteName:     "",
	//		Badge:         "",
	//		LocalFilePath: "",
	//	},
	//	{
	//		MenuUrl:    "/other/print",
	//		MenuName:   "打印",
	//		RouteName:  "",
	//		Icon:       "icon-printer",
	//		Cacheable:  true,
	//		ParentPath: "/other",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/other/clipboard",
	//		MenuName:   "剪贴板",
	//		RouteName:  "",
	//		Icon:       "icon-printer",
	//		Cacheable:  true,
	//		ParentPath: "/other",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/other/player",
	//		MenuName:   "视频播放器",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/other",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/other/flow",
	//		MenuName:   "流程图",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/other",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/excel",
	//		MenuName:   "Excel",
	//		RouteName:  "",
	//		Icon:       "icon-nav",
	//		Cacheable:  true,
	//		ParentPath: "",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/excel/export-excel",
	//		MenuName:   "导出Excel",
	//		RouteName:  "",
	//		Icon:       "icon-nav",
	//		Cacheable:  true,
	//		ParentPath: "/excel",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/excel/export-rows-excel",
	//		MenuName:   "导出选中行",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/excel",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/draggable",
	//		MenuName:   "拖拽",
	//		RouteName:  "",
	//		Icon:       "icon-drag-arrow",
	//		Cacheable:  true,
	//		ParentPath: "",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/draggable/card-draggable",
	//		MenuName:   "卡片拖拽",
	//		RouteName:  "",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/draggable",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/log",
	//		MenuName:   "系统监控",
	//		Icon:       "icon-computer",
	//		ParentPath: "",
	//		RouteName:  "log",
	//		Cacheable:  true,
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/log/operation",
	//		MenuName:   "操作日志",
	//		RouteName:  "operation",
	//		Icon:       "",
	//		Cacheable:  true,
	//		ParentPath: "/log",
	//		Badge:      "",
	//	}, {
	//		MenuUrl:    "/log/system",
	//		MenuName:   "系统日志",
	//		RouteName:  "system",
	//		Icon:       "icon-common",
	//		Cacheable:  true,
	//		ParentPath: "/log",
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/salary",
	//		MenuName:   "财务数据",
	//		Icon:       "icon-file",
	//		ParentPath: "",
	//		RouteName:  "salary",
	//		Cacheable:  true,
	//		Badge:      "",
	//	},
	//	{
	//		MenuUrl:    "/salary/salary",
	//		MenuName:   "工资详情",
	//		RouteName:  "salary",
	//		Icon:       "icon-empty",
	//		Cacheable:  true,
	//		ParentPath: "/salary",
	//		Badge:      "",
	//	},
	//}).Error
	//global.Global.Log.Error(err)
	//err = global.Global.Mysql.AutoMigrate(&models.Kind{})
	//err = global.Global.Mysql.AutoMigrate(&models.Tally{})
	//err = global.Global.Mysql.AutoMigrate(&models.Comment{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//err = global.Global.Mysql.AutoMigrate(&models.Collect{})
	//if err != nil {
	//	global.Global.Log.Info(err)
	//	return
	//}
	//err = global.Global.Mysql.AutoMigrate(&models.Comment{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//f := sync.Once{}
	//f.Do(
	//	func() {
	//		err = global.Global.Mysql.AutoMigrate(&models.User{})
	//		err = global.Global.Mysql.AutoMigrate(&models.Tally{})
	//		err = global.Global.Mysql.AutoMigrate(&models.Kind{})
	//	})
	//if err != nil {
	//	panic(err)
	//}
}
