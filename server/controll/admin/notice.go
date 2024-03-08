package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strconv"
	"time"
)

// PublishNotice 发布公告
func PublishNotice(c *gin.Context) {
	notice := new(global.Notice)
	err := c.Bind(notice)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	uid := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	var uids int
	if uid != "" {
		uids, err = strconv.Atoi(uid)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.AtoiError)
			return
		}
	} else {
		employer, err := dao.GetUserById(id)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.ClockInError)
			return
		}
		uids = int(int32(employer.Uid))
	}
	//
	var status int32 = 1
	date := time.Now().Unix()
	notices := &models.Notice{
		Identity: utils.GetUidV4(),
		Uid:      int64(uids),
		Title:    notice.Title,
		Text:     notice.Text,
		Status:   &status,
		Url:      notice.Url,
		Date:     date,
	}
	err = dao.InsertNotice(notices)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.PushNoticeError)
		return
	}
	//	放进缓存中
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, id, uids).Result()
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		//	放进zset结构
		marshal, err := json.Marshal(notices)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		global.Global.Redis.ZAdd(global.Global.Ctx, global.Notices, redis.Z{
			Score:  float64(date),
			Member: marshal,
		})
	})
	if err != nil {
		global.Global.Log.Error("goroutine  err:", err)
	}
	result.Ok(c, nil)
}

// UpdateNoticeStatus  修改公告状态
func UpdateNoticeStatus(c *gin.Context) {
	notice := new(global.UpdateNotice)
	err := c.Bind(notice)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryNotFoundError)
		return
	}
	//修改缓存
	_, err = global.Global.Redis.Del(global.Global.Ctx, global.Notices).Result()
	if err != nil {
		result.Fail(c, global.ServerError, global.UpdateNoticeError)
		return
	}
	err = dao.UpdateNotice(notice)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.UpdateNoticeError)
		return
	}
	result.Ok(c, nil)
}

func GetNoticeList(c *gin.Context) {
	//	先读缓存
	list, err := dao.GetNoticeLists()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.GetNoticeError)
		return
	}

	result.Ok(c, list)
}
