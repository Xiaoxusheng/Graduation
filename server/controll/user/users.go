package user

import (
	"github.com/gin-gonic/gin"
	"server/config"
	"server/dao"
	"server/global"
	"server/result"
	"server/utils"
)

type users struct {
	Password string `json:"password" form:"password" binding:"required" `
	Username string `json:"username" form:"username" binding:"required"`
}

// Login 登录
func Login(c *gin.Context) {
	user := new(users)
	err := c.Bind(user)
	if err != nil {
		result.Fail(c, global.UserCode, global.QueryError)
		return
	}
	global.Global.Log.Info(err)
	//胭脂

	//	判断用户是否存在
	//  查询盐值
	salt := global.Global.Redis.HGet(global.Global.Ctx, user.Username, "salt").Val()
	if salt == "" {
		result.Fail(c, global.UserCode, global.UserNotExist)
		return
	}
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Username, utils.HashPassword(user.Password, []byte(salt))).Val()
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
		global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time)
		result.Ok(c, map[string]any{
			"token": token,
		})
		return
	}
	//	identity 不存在
	use, err := dao.GetUserByUsePwd(user.Username, utils.HashPassword(user.Password, []byte(salt)))
	if err != nil || use == nil {
		result.Fail(c, global.UserCode, global.UserNotExist)
		return
	}
	//生成token
	token := utils.GetToken(use.Identity)
	global.Global.Redis.Set(global.Global.Ctx, val, token, config.Config.Jwt.Time)
	global.Global.Redis.HSet(global.Global.Ctx, user.Username, utils.HashPassword(user.Password, []byte(salt)))
	result.Ok(c, map[string]any{"token": token})
}

// Register 注册
func Register(c *gin.Context) {

}

// Info 用户信息
func Info(c *gin.Context) {

}

//注销

//退出登录
