package config

import (
	"server/global"
	"sync"
)

func InitMutex() {
	global.Global.Mutex = &sync.RWMutex{}
}
