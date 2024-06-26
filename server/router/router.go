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
	//菜单列表
	e.POST("/admin/menu_list", middleware.ParseToken(), menu.GetMenuList)
	//
	//u := e.Group("/user")
	//e.POST("/login", user.Login)

	roots := e.Group("/root")
	roots.Use(middleware.Log(), middleware.ParseToken(), middleware.CasBin())
	//分配角色
	roots.POST("/add_rolesForUser", root.AddRolesForUser)
	//分配资源
	roots.POST("/add_permissionForUser", root.AddPermissionForUser)
	//删除用户的角色
	roots.POST("/delete_roleForUser", root.DeleteRoleForUser)
	//删除角色,包括删除角色所有权限
	roots.POST("/delete_role", root.DeleteRole)
	//删除单个权限
	roots.POST("/delete_permissionForUser", root.DeletePermissionForUser)
	//删除用户所有的权限
	roots.POST("/delete_Permission", root.DeletePermissionsForUser)
	//更新权限
	roots.POST("/update_policy", root.UpdatePolicy)
	//查看角色的权限
	roots.POST("/get_permissionsForUser", root.GetPermissionsForUser)
	//查看所有权限
	roots.POST("/get_allNamedSubjects", root.GetAllNamedSubjects)
	//获取不同角色的菜单列表
	roots.GET("/get_role_menuList", root.GetRoleMenuList)
	//添加角色能访问的菜单
	roots.POST("/add_role_menu", root.AddRoleMenu)
	//删除角色能访问的菜单
	roots.POST("/del_role_menu", root.DelRoleMenu)
	//更新菜单
	roots.POST("/update_menu", root.UpdateMenu)
	//获取所有的角色
	roots.GET("/role_list", root.GetRoleList)

	roots.POST("/add", root.Add)
	//登录
	e.POST("/admin/login", admin.Login)

	api := e.Group("/admin")
	api.Use(middleware.Log(), middleware.ParseToken(), middleware.CasBin())
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
	//获取不同状态一天的考勤
	api.GET("/get_by_status", admin.GetByStatus)
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

	//工资
	//查员工的工资
	api.GET("/get_salary", admin.GetSalary)
	//获取所有的工资列表
	api.GET("/get_salary_list", admin.GetSalaryList)
	//输入信息
	api.POST("/salary_info", admin.SalaryInfo)
	//删除工资信息
	api.GET("/delete_salary", admin.DeleteSalary)

	//增加菜单
	api.POST("/add_menu", menu.AddMenu)
	//删除菜单
	api.GET("/del_menu", menu.DeleteMenu)
	//更新菜单
	api.POST("/update_menu", menu.UpdateMenu)

	//发布公告
	api.POST("/publish_notice", admin.PublishNotice)
	//更新公告
	api.POST("/update_notice", admin.UpdateNoticeStatus)
	//获取公告列表
	api.GET("/get_notice_list", admin.GetNoticeList)
	//删除公告
	api.GET("/del_notice", admin.DelNotice)

	//日志
	api.GET("/log_list", admin.GetLogList)
	api.GET("/user_log_list", admin.GetUserLogList)

	//------------------------------------------------------
	users := e.Group("/user", middleware.Log(), middleware.ParseToken())
	//文件上传
	users.POST("/upload", user.Upload)
	//打卡
	users.GET("/clockIn", user.ClockIn)
	//补卡申请
	users.POST("/markCard_application", user.MarkCardApplication)
	//请假申请
	users.POST("/leave", user.LeaveApplication)
	//加班申请
	users.POST("/overtime", user.OverTimeApplication)
	//获取个人信息
	users.GET("/info", user.EmployeeInfo)
	//修改个人信息

	//查看考勤列表
	users.GET("/get_clockIn_log", user.GetClockInLog)
	//获取个人工资情况
	users.GET("/get_salary", user.GetSalary)
	//查看申请
	users.GET("/get_examine", user.GetExamine)

	//获取公告
	users.GET("/get_notice_list", user.GetNotice)
	//修改信息
	users.POST("/change_info", user.ChangeInfo)

	//修改密码
	users.GET("/change_password", user.ChangePassword)
	//
	return e
}
