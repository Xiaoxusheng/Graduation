package global

// *错误码*/
// 通用错误
const (
	UnknownError       = 20001 // 未知错误
	ServiceUnavailable = 20002 // 服务不可用
)

// 客户端错误
const (
	BadRequest           = 20010 // 请求参数错误
	AuthenticationFailed = 20011 // 身份验证失败
	PermissionDenied     = 20012 // 权限不足
	ResourceNotFound     = 20013 // 资源不存在
	MethodNotAllowed     = 20014 // 方法不支持
)

// 服务器错误
const (
	ServerError            = 30001 // 服务器内部错误
	NotImplemented         = 30002 // 功能未实现
	ServiceTimeout         = 30003 // 服务超时
	DependencyServiceError = 30004 // 依赖服务失败
	DatabaseError          = 30005 // 数据库错误
	ExceedLimitError       = 30006 // 服务器请求超过限制
)

// 数据错误
const (
	DataValidationError = 40001 // 数据验证失败
	DataConflict        = 40002 // 数据冲突
	DataNotFound        = 40003 // 数据不存在
	DataTooLarge        = 40004 // 数据过大
	DataUnmarshal       = 40005 // 数据解析失败
)

// *错误信息*/
var (
	TokenError             = "token格式错误"
	FileError              = "文件数量超过限制！"
	QueryError             = "参数错误！"
	PhoneError             = "手机号已经存在！"
	DataUnmarshalError     = "数据解析失败！"
	UserNotExistError      = "用户不存在或密码错误！"
	QueryNotFoundError     = "获取必要参数失败！"
	ResetPwdError          = "重置密码失败"
	ChangePwdError         = "修改密码失败！"
	ParseError             = "解析错误！"
	RequestToError         = "请求无法处理！"
	EmployerNotFoundError  = "员工不存在！"
	UserNotFound           = "用户不存在！"
	DeleteError            = "删除失败！"
	AddEmployerError       = "新增员工失败！"
	UpdateEmployerError    = "更新员工信息失败！"
	AssignedAccountError   = "分配账号失败！"
	GetEmployerListError   = "获取员工列表失败！"
	GetDepartmentError     = "获取部门列表失败！"
	GetDepartmentInfoError = "获取部门信息失败!"
	DelDepartmentError     = "删除部门失败！"
	DepartmentNotFound     = "部门不存在！"
	UpdateDepartmentError  = "更新部门失败！"
	AddDepartmentError     = "新增部门失败！"
	UpdateClockInLogError  = "修改考勤记录失败！"

	//菜单
	AddMenuError    = "新增菜单失败！"
	DeleteMenuError = "删除菜单失败！"
	UpdateMenuError = "更新菜单失败！"

	//考勤
	GetClockError = "获取考勤列表失败！"
	AtoiError     = "转化失败！"

	//	请假列表
	GetLeaveListError     = "获取请假列表失败！"
	LeaveApplicationError = "请假审批失败！"

	//	加班
	OverTimeApplicationError = "加班申请审批失败！"
	GetOverTimeError         = "获取加班申请列表失败！"

	//	补卡
	MarkCardApplicationError = "补卡申请审核失败！"
	GetMarkCardLiatError     = "获取补卡申请列表失败！"

	//	日志
	GetLogList = "获取日志列表失败！"

	//	工资
	GetSalaryError     = "获取工资失败！"
	GetSalaryListError = "获取工资列表失败！"
)

var (
	UserRepeatClockIn = "重复打卡"
	ClockInError      = "打卡失败！"
	MarkCardError     = "补卡失败！"
	LeaveError        = "请假失败！"
	OverTimeError     = "加班申请失败！"
)

// 超级管理员
var (
	PermissionDeniedError     = "权限不足！"
	AddRoleFail               = "添加角色失败！"
	AddPermissionFail         = "分配资源失败！"
	RoleNotfound              = "角色不存在！"
	DelRoleFail               = "删除角色失败！"
	DelPermissionFail         = "删除资源失败！"
	UpdatePermissionFail      = "更新资源失败！"
	PermissionNotFound        = "资源不存存在！"
	GetPermissionsForUserFail = "获取列表失败！"
)
