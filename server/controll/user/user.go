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
			result.Ok(c, map[string]any{
				"token": token,
			})
			return
		}
		//   token不存在
		token = utils.GetToken(val)
		global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time*time.Minute*60)
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
		global.Global.Redis.Set(global.Global.Ctx, use.Identity, token, config.Config.Jwt.Time)
		global.Global.Redis.HSet(global.Global.Ctx, user.Account, utils.HashPassword(user.Password, salts), use.Identity)
	}()
	result.Ok(c, map[string]any{"token": token})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {

}
