package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"strconv"
	"time"
)

// GetNotice 获取公告
func GetNotice(c *gin.Context) {
	//	先读缓存
	list := make([]*models.Notice, 10)
	val := global.Global.Redis.ZRangeByScoreWithScores(global.Global.Ctx, global.Notices, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(time.Now().Unix(), 10),
		Offset: 0,
		Count:  10,
	}).Val()
	if len(val) != 0 {
		for i := 0; i < len(val); i++ {
			notice := new(models.Notice)
			err := json.Unmarshal([]byte(val[i].Member.(string)), notice)
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.GetNoticeError)
				return
			}
			list = append(list, notice)
		}
		result.Ok(c, list)
		return
	}
	list, err := dao.GetNoticeList()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.GetNoticeError)
		return
	}
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		for i := 0; i < len(list); i++ {
			marshal, err := json.Marshal(list[i])
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
			_, err = global.Global.Redis.ZAdd(global.Global.Ctx, global.Notices, redis.Z{
				Score:  float64(list[i].Date),
				Member: marshal,
			}).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
		}
	})
	if err != nil {
		global.Global.Log.Error(err)
	}
	result.Ok(c, list)
}
