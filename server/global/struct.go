package global

import (
	"time"
)

// UrlList 图片返回结构体
type UrlList struct {
	Url   string `json:"url,omitempty"`
	Index int    `json:"index,omitempty"`
}

// Department 部门
type Department struct {
	Identity string `json:"identity,omitempty"`
	Name     string `json:"name,omitempty"  binding:"required,min=2,max=15" form:"name" form:"name"`
	Sort     int32  `json:"sort,omitempty" binding:"required,number" form:"sort" form:"sort"`
	Leader   string `json:"leader,omitempty" binding:"required,min=2,max=15" form:"leader" form:"leader"`
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
	Uid        int64  `json:"uid,omitempty" binding:"required" form:"uid"`
	Name       string `json:"name,omitempty"  binding:"required" form:"name"`
	Birthday   int64  `json:"birthday,omitempty" binding:"required" form:"birthday"`
	Status     int32  `json:"status" binding:"required" form:"status"`
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
	Uid       int64  `json:"uid,omitempty"  binding:"required" form:"uid"`
	StartTime int64  `json:"start_time,omitempty"  binding:"required" form:"startTime"`
	EndTime   int64  `json:"end_time,omitempty"  binding:"required" form:"endTime"`
	Pass      int32  `json:"pass"  binding:"required" form:"pass"`
	Reason    string `json:"reason,omitempty"  binding:"required" form:"reason"`
	Limit     int32  `json:"limit" binding:"required" form:"limit"`
	Offset    int32  `json:"offset" binding:"required" form:"offset"`
}
