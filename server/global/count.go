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
	Sign           = "sign"           //登录
	Salt           = "salt"           //盐值
	Info           = "info"           //个人信息
	Uid            = "uid"            //uid
	Employer       = "employer"       //员工
	DepartmentList = "departmentList" //部门
	DepartmentId   = "department_id"  //存放部门identity
	GetClockInLog  = "GetClockInLog"  //考勤记录
)

// redis 时间
const (
	InfoTime         = 60 * 60 * 24 * 7
	DepartmentTime   = 60 * 60 * 24 * 7
	EmployerInfoTime = 60 * 60 * 24
)
