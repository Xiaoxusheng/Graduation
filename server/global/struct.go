package global

type Department struct {
	Identity string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity,omitempty" binding:"required min=2 max=10"`
	Name     string `gorm:"type:varchar(10);not null unique; comment:'部门名称'" json:"name,omitempty" binding:"required min=2 max=15"`
	Sort     int32  `gorm:"type:int not null unique; comment:'部门编号'" json:"sort,omitempty" binding:"required number"`
	Leader   string `gorm:"type:varchar(36);not null unique; comment:'主管'" json:"leader,omitempty" binding:"required min=2 max=15"`
}
