package router

import (
	"github.com/gin-gonic/gin"
	"server/controll/admin"
	"server/controll/menu"
	"server/controll/root"
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

	roots := e.Group("/root")
	roots.Use(middleware.ParseToken(), middleware.CasBin())
	//分配角色
	roots.POST("/add_rolesForUser", root.AddRolesForUser)
	//分配资源
	roots.POST("/add_resource", root.AddResource)
	//删除角色
	roots.POST("/delete_roleForUser", root.DeleteRoleForUser)
	//删除资源
	roots.POST("/delete_permissionForUser", root.DeletePermissionForUser)
	//更新权限
	roots.POST("/update_policy", root.UpdatePolicy)
	//查看能访问的资源
	roots.POST("/get_permissionsForUser", root.GetPermissionsForUser)
	//查看所有权限
	roots.POST("/get_allNamedSubjects", root.GetAllNamedSubjects)

	//登录
	e.POST("/admin/login", admin.Login)

	api := e.Group("/admin")
	api.Use(middleware.ParseToken(), middleware.CasBin())
	//注册
	e.POST("/register", admin.Register)
	//个人信息
	api.GET("/info", admin.Info)
	//退出登录
	api.GET("/logout", admin.Logout)
	//重置密码
	api.GET("/reset_password", admin.ResetPassword)

	//删除员工
	api.GET("/delete_employer", admin.DeleteEmployee)
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
	//获取部门人数信息
	api.GET("/department_info", admin.DepartmentInfo)

	//获取员工考勤记录
	api.GET("/get_clockIn", admin.GetClockInLog)
	//编辑考勤记录
	api.POST("/edit_clockIn", admin.EditClockLog)
	//获取某一天全部考勤
	api.GET("/get_all_clockIn", admin.GetClockList)

	//请假申请审核
	api.POST("/leave_application", admin.LeaveApplication)
	//获取请假申请列表
	api.GET("/get_LeaveApplication_list", admin.GetLeaveApplicationList)
	//加班申请审核
	api.POST("/overtime_application", admin.OvertimeApplication)
	//获取加班申请表
	api.GET("/get_overtimeApplication_list", admin.GetOvertimeList)
	//补卡申请审批
	api.POST("/make_card_application", admin.MakeCardApplication)
	//补卡申请列表
	api.GET("/make_card_application_list", admin.GetMarkCardList)

	//增加菜单
	api.POST("/add_menu", menu.AddMenu)

	//------------------------------------------------------
	users := e.Group("/user", middleware.ParseToken())
	//文件上传
	users.POST("/upload", user.Upload)
	//打卡
	users.GET("/clockIn", user.ClockIn)
	//补卡申请
	users.POST("/markCard_application", user.MarkCardApplication)
	//请假
	users.POST("/leave", user.LeaveApplication)
	//加班
	users.POST("/overtime", user.OverTimeApplication)

	//修改密码
	users.GET("/change_password", user.ChangePassword)
	//
	return e
}
