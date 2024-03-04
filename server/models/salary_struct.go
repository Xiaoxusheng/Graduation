package models

// SalaryStruct 薪资结构表
type SalaryStruct struct {
	Identity   string `gorm:"type:varchar(36) not null unique; comment:'工资记录唯一标识'" json:"identity,omitempty"`
	Position   int32  `gorm:"type:int not null; comment:'职位'" json:"Position,omitempty"`
	BaseSalary int32  `gorm:"type:int not null; comment:''" json:"base_salary,omitempty"`
	Allowance  int32  `gorm:"type:int not null; comment:'补贴'" json:"allowance,omitempty"`
	Other      int32  `gorm:"type:int not null; comment:'其他扣除'" json:"other,omitempty"`
}
