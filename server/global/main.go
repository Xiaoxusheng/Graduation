package global

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Configs struct {
	Mysql *gorm.DB        `json:"mysql"`
	Redis *redis.Client   `json:"redis"`
	Log   *log.Logger     `json:"log"`
	Ctx   context.Context `json:"ctx"`
}

var (
	Global Configs
)
