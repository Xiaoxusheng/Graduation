package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "xlei.love:6379",
		Password: "admin123", // 没有密码，默认值
		DB:       0,
	}).WithTimeout(time.Second * 5)
	ctx := context.Background()
	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if res == "PONG" {
		log.Println("redis连接成功!")
	}
}
