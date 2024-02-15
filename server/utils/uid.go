package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"server/global"
	"time"
)

func GetUidV5(u string) string {
	return uuid.NewV5(uuid.NewV4(), u).String()
}

func GetUidV4() string {
	return uuid.NewV4().String()
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// CreateUid 生成唯一id
func CreateUid(key string) int64 {
	t := time.Now()
	t1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	//获取redis中的数据
	val := global.Global.Redis.IncrBy(global.Global.Ctx, key+t.Format(time.DateOnly), 1).Val()
	return t.Unix() - t1.Unix()<<32 | val
}
