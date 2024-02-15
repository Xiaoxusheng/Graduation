package middleware

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

// TokenBucket 桶
type TokenBucket struct {
	mutex    sync.RWMutex  //读写锁
	capacity int32         //令牌桶容量
	rate     time.Duration //生成速率
	current  int32         //当前令牌桶令牌数量
}

func NewTokenBucket(capacity int32, rate int64) *TokenBucket {
	return &TokenBucket{
		mutex:    sync.RWMutex{},
		capacity: capacity,
		rate:     time.Second * time.Duration(rate),
		current:  capacity,
	}
}

func (t *TokenBucket) limit() bool {

	return true

}

// RateLimit 限速中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
