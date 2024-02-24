package router

import (
	"github.com/gin-gonic/gin"
	"server/controll/admin"
	"server/controll/menu"
	"server/controll/user"
	"server/middleware"
)

func Routers(e *gin.Engine) *gin.Engine {
	//

	e.POST("/user/login", user.Login)
	e.POST("/admin/menu_list", menu.GetMenuList)
	//
	//u := e.Group("/user")
	//e.POST("/login", user.Login)

	//登录
	e.POST("/admin/login", admin.Login)

	api := e.Group("/admin")
	api.Use(middleware.ParseToken())
	//注册
	e.POST("/register", admin.Register)
	//个人信息
	api.GET("/info", admin.Info)
	//退出登录
	api.GET("/Logout", admin.Logout)
	//重置密码
	api.GET("reset_password", admin.ResetPassword)

	//删除员工
	api.GET("delete_employer", admin.DeleteEmployee)
	//创建员工信息
	api.POST("/add_employer", admin.AddEmployee)
	//更新员工信息
	api.POST("/update_employer", admin.UpdateEmployee)
	//员工信息列表
	api.GET("/employer_list", admin.EmployeeList)
	//员工个人信息
	api.GET("/employer_info", admin.EmployeeInfo)

	//分配账号
	api.GET("/assigned_account", admin.AssignedAccount)

	//部门列表
	api.GET("/get_department_list", admin.GetDepartmentList)
	//删除部门
	api.GET("/del_department", admin.DeleteDepartment)
	//更新部门信息
	api.POST("/update_department", admin.UpdateDepartment)
	//新增部门信息
	api.POST("/add_department", admin.AddDepartment)

	//获取考勤记录
	api.GET("/get_clockIn", admin.GetClockInLog)
	//编辑考勤记录
	api.POST("edit_clockIn", admin.EditClockLog)

	//请假申请审核
	api.POST("/leave_application", admin.LeaveApplication)
	//获取请假申请列表
	api.GET("/get_LeaveApplication_list", admin.GetLeaveApplicationList)
	//加班申请审核
	api.POST("overtime_application", admin.OvertimeApplication)
	//获取加班申请表
	api.GET("/get_overtimeApplication_list", admin.GetOvertimeList)
	//补卡申请审批
	api.POST("/make_card_application", admin.MakeCardApplication)
	//补卡申请列表
	api.POST("/make_card_application_liat", admin.GetMarkCardList)

	//------------------------------------------------------
	e.Group("/user", middleware.ParseToken())
	//文件上传
	e.POST("/upload", user.Upload)

	return e
}
