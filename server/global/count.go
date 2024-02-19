package global

// 日=日志颜色
const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Gray   = 38
)

// redis 前缀
var (
	Sign           = "sign"
	Salt           = "salt"
	UidKey         = "uid-key"
	Info           = "info"
	Uid            = "uid"
	Employer       = "employer"
	DepartmentList = "departmentList"
)

// redis 时间
const (
	InfoTime       = 60 * 60 * 24 * 7
	DepartmentTime = 60 * 60 * 24 * 7
)
