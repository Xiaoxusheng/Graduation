package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
	"sync"
	"sync/atomic"
	"time"
)

// TokenBucket 桶
type TokenBucket struct {
	mutex    sync.RWMutex
	now      time.Time //当前时间
	capacity int32     //令牌桶容量
	rate     int64     //生成速率
	current  int32     //当前令牌桶令牌数量
}

func NewTokenBucket(capacity int32, rate int64) *TokenBucket {
	return &TokenBucket{
		mutex:    sync.RWMutex{},
		now:      time.Now(),
		capacity: capacity,
		rate:     rate,
		current:  capacity,
	}
}

func (t *TokenBucket) limit() bool {
	/*原子操作*/
	//获取当前时间
	now := time.Now()
	t1 := now.UnixMilli() - t.now.UnixMilli()
	fmt.Println(now.UnixMilli(), t.now.UnixMilli())
	fmt.Println("时间", t1)

	num := t1 / 200 * t.rate
	fmt.Println("生成的令牌", num)
	//生成的token大于总容量
	if t.current+int32(num) >= t.capacity {
		atomic.StoreInt32(&t.current, t.capacity)
	} else {
		atomic.StoreInt32(&t.current, t.current+t.capacity)
	}
	//
	current := atomic.LoadInt32(&t.current)
	if current > 0 {
		//令牌减一
		atomic.AddInt32(&t.current, -1)
		fmt.Println("令牌数量", current)
		//    重置时间
		t.mutex.Lock()
		defer t.mutex.Unlock()
		t.now = time.Now()
		return true
	}
	return false

}

// RateLimit 限速中间件
func RateLimit() gin.HandlerFunc {
	tokenBucket := NewTokenBucket(20, 1)
	return func(c *gin.Context) {
		ok := tokenBucket.limit()
		if !ok {
			result.Fail(c, global.ExceedLimitError, global.RequestToError)
			c.Abort()
			return
		}
		c.Next()
	}
}
