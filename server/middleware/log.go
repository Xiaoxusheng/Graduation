package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/utils"
	"strconv"
	"time"
)

// Log 记录日志中间件
func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	获取请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		t := time.Now()
		c.Next()
		httpCode := c.Writer.Status()
		id := c.GetString("identity")
		val := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
		if val == "" {
			employer, err := dao.GetUid(id)
			if err != nil || employer.Account == 0 {
				global.Global.Log.Error(err)
				return
			}
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, id, employer.Account).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
		}
		err := global.Global.Pool.Submit(func() {
			val = global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
			atoi, err := strconv.Atoi(val)
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
			//	写入数据库
			err = dao.InsertLog(&models.Log{
				Identity: utils.GetUidV4(),
				Method:   method,
				Path:     path,
				IP:       c.RemoteIP(),
				Time:     time.Now().Sub(t).Milliseconds(),
				Uid:      int64(atoi),
				HttpCode: int32(httpCode),
			})
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
		})
		if err != nil {
			global.Global.Log.Error(err)
		}
		global.Global.Log.Info(time.Now().Sub(t))
		fmt.Println(path, method, id, httpCode)
		if path == "" || method == "" || id == "" || httpCode == 0 {
			return
		}
		fmt.Println(path, method, id, httpCode, time.Now().Sub(t).Milliseconds())

	}
}
