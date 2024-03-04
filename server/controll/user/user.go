package user

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"server/config"
	"server/dao"
	"server/global"
	"server/result"
	"server/utils"
	"time"
)

type users struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required" `
}

// Login 登录
func Login(c *gin.Context) {
	user := new(users)
	err := c.Bind(user)
	if err != nil {
		global.Global.Log.Info(err)
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//	判断用户是否存在
	//  查询盐值
	salt := global.Global.Redis.HGet(global.Global.Ctx, user.Account, global.Salt).Val()
	//salt := global.Global.Redis.HGet(global.Global.Ctx, global.UidKey+user.Account, global.Salt).Val()
	if salt == "" {
		//数据库获取
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.UserNotExistError)
		return
	}
	salts, _ := base64.URLEncoding.DecodeString(salt)
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Account, utils.HashPassword(user.Password, salts)).Val()
	//identity 存在
	if val != "" {
		//获取token
		token := global.Global.Redis.Get(global.Global.Ctx, val).Val()
		if token != "" {
			//判断身份，返回身份信息
			result.Ok(c, map[string]any{
				"token": token,
			})
			return
		}
		//   token不存在
		token = utils.GetToken(val)
		global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time*time.Hour)
		result.Ok(c, map[string]any{
			"token": token,
		})
		return
	}
	//	identity 不存在
	use, err := dao.GetUserAccountPwd(user.Account, utils.HashPassword(user.Password, salts))
	if err != nil || use == nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.UserNotExistError)
		return
	}
	//生成token
	token := utils.GetToken(use.Identity)
	go func() {
		global.Global.Redis.Set(global.Global.Ctx, use.Identity, token, config.Config.Jwt.Time*time.Hour)
		global.Global.Redis.HSet(global.Global.Ctx, user.Account, utils.HashPassword(user.Password, salts), use.Identity)
	}()
	result.Ok(c, map[string]any{"token": token})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	pwd := c.Query("password")
	p := c.Query("pwd")
	if p == "" || pwd == "" {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//查询uid
	uid := global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	//查询盐值
	val := global.Global.Redis.HGet(global.Global.Ctx, uid, global.Salt).Val()
	if val != "" {
		salt, _ := base64.URLEncoding.DecodeString(val)
		err := dao.UpdatePwd(uid, utils.HashPassword(pwd, salt))
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ChangePwdError)
			return
		}
		//
		_, err = global.Global.Redis.HDel(global.Global.Ctx, uid, utils.HashPassword(p, salt)).Result()
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ChangePwdError)
			return
		}
		_, err = global.Global.Redis.HSet(global.Global.Ctx, uid, utils.HashPassword(pwd, salt), id).Result()
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ChangePwdError)
			return
		}
		result.Ok(c, nil)
		return
	}
	//盐值不存在，通过工号查数据库查所有信息
	account, err := dao.GetByAccount(uid)
	if err != nil || account == nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.UserNotFound)
		return
	}
	//盐值
	salt, _ := base64.URLEncoding.DecodeString(account.Salt)
	//更新
	err = dao.UpdatePwd(uid, utils.HashPassword(pwd, salt))
	if err != nil {
		global.Global.Log.Warn(err)
		result.Fail(c, global.BadRequest, global.ResetPwdError)
		return
	}
	_, err = global.Global.Redis.HDel(global.Global.Ctx, uid, utils.HashPassword(p, salt)).Result()
	if err != nil {
		global.Global.Log.Warn(err)
		result.Fail(c, global.BadRequest, global.ChangePwdError)
		return
	}
	_, err = global.Global.Redis.HSet(global.Global.Ctx, uid, utils.HashPassword(pwd, salt), id).Result()
	if err != nil {
		global.Global.Log.Warn(err)
		result.Fail(c, global.BadRequest, global.ChangePwdError)
		return
	}
	result.Ok(c, nil)

}
