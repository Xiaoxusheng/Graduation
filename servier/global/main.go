package global

import (
	"context"
	"gorm.io/gorm"
	"log"
)

type Configs struct {
	Mysql *gorm.DB `json:"mysql"`
	//Redis  *redis.Client          `json:"redis"`
	Log *log.Logger     `json:"log"`
	Ctx context.Context `json:"ctx"`
}

var (
	Global Configs
)
