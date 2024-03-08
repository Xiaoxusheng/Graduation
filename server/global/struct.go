package global

import (
	"time"
)

// UrlList 图片返回结构体
type UrlList struct {
	Url   string `json:"url,omitempty"`
	Index int    `json:"index,omitempty"`
}

type AdminInfo struct {
	Sex          int32     `json:"sex"`
	Birthday     time.Time `json:"birthday"`
	DepartmentId int32     `json:"department_id"`
	Position     int32     `json:"position"`
	Uid          int64     `json:"uid"`
	NickName     string    `json:"nickName"`
	ImageUrl     string    `json:"avatar"`
	IP           string    `json:"IP"`
	Identity     string    `json:"user_identity,"`
	Username     string    `json:"username"`
}

// Department 部门
type Department struct {
	Identity string `json:"identity,omitempty" form:"identity"`
	Name     string `json:"name,omitempty"  binding:"required,min=2,max=15" form:"name" `
	Sort     int32  `json:"sort,omitempty" binding:"required,number" form:"sort" `
	Status   int32  `json:"status,omitempty" binding:"required,number" form:"status" `
	Leader   string `json:"leader,omitempty" binding:"required,min=2,max=15" form:"leader" `
}

// UserInfo 个人信息
type UserInfo struct {
	Identity     string    `json:"identity,omitempty"`
	Id           int32     `json:"id"`
	Uid          int64     `json:"uid,omitempty"`
	Name         string    `json:"name,omitempty"`
	Birthday     time.Time `json:"birthday,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	ImageUrl     string    `json:"image_Url"`
	IP           string    `json:"IP,omitempty"`
	Status       int32     `json:"status"`
	Sex          int32     `json:"sex"`
	DepartmentId int32     `json:"department_id,omitempty"`
	Position     int32     `json:"position,omitempty"`
}

// Employers 增加员工信息
type Employers struct {
	Uid        int64  `json:"uid,omitempty"  form:"uid"`
	Name       string `json:"name,omitempty"  binding:"required" form:"name"`
	Birthday   int64  `json:"birthday,omitempty" binding:"required" form:"birthday"`
	Status     int32  `json:"status" form:"status"`
	Sex        int32  `json:"sex" binding:"required" form:"sex"`
	Phone      string `json:"phone,omitempty" binding:"required" form:"phone"`
	Position   int32  `json:"position,omitempty" binding:"required,number" form:"position"`
	Department int32  `json:"department,omitempty" binding:"required,number" form:"department"`
}

// Attendance 考勤信息
type Attendance struct {
	Identity  string `json:"identity,omitempty"  binding:"required"`
	Username  string `json:"username,omitempty"  binding:"required"`
	Uid       int64  `json:"uid,omitempty"  binding:"required"`
	Name      string `json:"name,omitempty"  binding:"required"`
	StartTime int64  `json:"start_time,omitempty"  binding:"required"`
	EndTime   int64  `json:"end_time,omitempty"  binding:"required"`
	Status    int32  `json:"status"  binding:"required"`
	Reason    string `json:"reason,omitempty"  binding:"required"`
	Url       string `json:"url"  binding:"required"`
}

// Application 申请信息
type Application struct {
	Uid  int64 `json:"uid,omitempty"  binding:"required" form:"uid"`
	Pass int32 `json:"pass"  binding:"required" form:"pass"`
}

// Menu 菜单
type Menu struct {
	Uid           int32  `json:"uid" form:"uid"`
	MenuUrl       string `json:"menuUrl,omitempty" form:"menuUrl"`
	MenuName      string `json:"menuName,omitempty" form:"menuName"`
	Icon          string `json:"icon,omitempty" form:"icon"`
	ParentPath    string `json:"parentPath,omitempty" form:"parentPath"`
	RouteName     string `json:"routeName,omitempty" form:"routeName"`
	Cacheable     bool   `json:"cacheable,omitempty" form:"cacheable"`
	Badge         string `json:"badge,omitempty" form:"badge"`
	LocalFilePath string `json:"localFilePath,omitempty" form:"localFilePath"`
	IsRootPath    bool   `json:"isRootPath" form:"isRootPath"`
	Hidden        bool   `json:"hidden" form:"hidden"`
}

// MarkCard 补卡申请
type MarkCard struct {
	Date   int64  `json:"date"    binding:"required"   form:"date"`
	Url    string `json:"url,omitempty"   binding:"required"  form:"url"`
	Reason string `json:"reason,omitempty"   binding:"required" form:"reason"`
}

// LeaveApplication 请假
type LeaveApplication struct {
	StartTime int64  `json:"start_time" binding:"required"   form:"start_time"`
	EndTime   int64  `json:"end_time"  binding:"required"   form:"end_time"`
	Reason    string `json:"reason"   binding:"required"  form:"reason"`
	Url       string `json:"url,omitempty"   binding:"required"  form:"url"`
}

// OverTime 加班
type OverTime struct {
	StartTime int64  `json:"start_time" binding:"required"   form:"start_time"`
	EndTime   int64  `json:"end_time"  binding:"required"   form:"end_time"`
	Url       string `json:"url,omitempty"   binding:"required"  form:"url"`
}

// Applications 申请
type Applications struct {
	Uid          int64     `json:"uid,omitempty"`
	Pass         int32     `json:"pass"`
	IsExamine    int32     `json:"is_examine"`
	Sex          int32     `json:"sex"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	Status       int32     `json:"status"`
	DepartmentId int32     `json:"department_id"`
	Reason       string    `json:"reason,omitempty"`
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	Identity     string    `json:"identity,omitempty"`
}

// ClockingIn 打卡
type ClockingIn struct {
	Identity     string    `json:"identity,omitempty"`
	Uid          int64     `json:"uid,omitempty"`
	Sex          int32     `json:"sex"`
	Name         string    `json:"name,omitempty"`
	DepartmentId int32     `json:"department_id,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	EndTime      time.Time `json:"end_time,omitempty"`
	Date         time.Time `json:"date,omitempty"`
	Status       int32     `json:"status"`
}

// Notice 公告
type Notice struct {
	Text  string `json:"text" form:"text"  binding:"required" `
	Title string `json:"title" form:"title"  binding:"required" `
	Url   string `json:"url" form:"url"  binding:"required" `
}

type UpdateNotice struct {
	Id     string `json:"id" form:"id"  binding:"required" `
	Status int32  `json:"status" form:"status"  binding:"required" `
	Text   string `json:"text" form:"text"  binding:"required" `
	Title  string `json:"title" form:"title"  binding:"required" `
	Url    string `json:"url" form:"url"  binding:"required" `
}
