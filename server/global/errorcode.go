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
	ServerError            = 30001 // 内部服务器错误
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
)

// *错误信息*/
var (
	QueryError            = "参数错误！"
	UserNotExistError     = "用户不存在！"
	QueryNotFoundError    = "获取必要参数失败！"
	ParseError            = "解析错误！"
	RequestToError        = "请求无法处理！"
	EmployerNotFoundError = "员工不存在！"
	DeleteError           = "删除失败！"
	AddEmployerError      = "新增员工失败！"
)
