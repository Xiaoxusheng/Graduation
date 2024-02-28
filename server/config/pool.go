package config

import (
	"github.com/panjf2000/ants/v2"
	"server/global"
	"sync"
)

func InitPool() {
	defer ants.Release()
	mp, err := ants.NewMultiPool(Config.Pool.Num, Config.Pool.Size, ants.RoundRobin)
	if err != nil {
		panic(err)
	}
	global.Global.Pool = mp
	global.Global.Wg = new(sync.WaitGroup)
}
