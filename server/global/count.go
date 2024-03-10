package global

// 日=日志颜色
const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Gray   = 38
)

const FileNumber = 4

// redis 前缀
var (
	Salt               = "salt"                 //盐值
	Info               = "info"                 //个人信息
	Uid                = "uid"                  //uid
	Employer           = "employer"             //员工
	DepartmentList     = "departmentList"       //部门
	DepartmentId       = "department_id"        //存放部门identity
	GetClockInLog      = "GetClockInLog"        //考勤记录
	Menus              = "menu"                 //菜单
	ClockIn            = "ClockIn"              //打卡
	UidId              = "uid-identity"         //员工的uid和identity的映射
	User               = "User"                 //超级管理员
	SalaryList         = "salary-list"          //工资列表
	SalaryEmployerList = "salary_employer_list" //员工个人工资列表
	Notices            = "notice"               //员工获取的公告
)

// redis 时间
const (
	InfoTime               = 60 * 60 * 24 * 7
	DepartmentTime         = 60 * 60 * 24 * 7
	EmployerClockTime      = 60 * 60
	MenuTime               = 60 * 60 * 24
	EmployerUidId          = 60 * 60 * 24 * 7
	SalaryListTime         = 60 * 60
	SalaryEmployerListTime = 60 * 60
)
