package global

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

type Configs struct {
	Mysql  *gorm.DB         `json:"mysql"`
	Redis  *redis.Client    `json:"redis"`
	Log    *log.Logger      `json:"log"`
	Ctx    context.Context  `json:"ctx"`
	Mutex  *sync.RWMutex    `json:"mutex"`
	Pool   *ants.MultiPool  `json:"pool"`
	Wg     *sync.WaitGroup  `json:"wg"`
	CasBin *casbin.Enforcer `json:"casBin"`
}

var (
	Global Configs
)
