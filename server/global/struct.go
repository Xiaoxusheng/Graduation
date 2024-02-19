package global

type Department struct {
	Name   string `json:"name,omitempty" binding:"required,min=2,max=15" form:"name"`
	Sort   int32  `json:"sort,omitempty" binding:"required,number" form:"sort"`
	Leader string ` json:"leader,omitempty" binding:"required,min=2,max=15" form:"leader"`
}
