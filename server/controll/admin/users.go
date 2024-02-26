package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/config"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strconv"
	"time"
)

type users struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required" `
}

type registers struct {
	Username string `json:"username"  form:"username"  validate:"min=5,max=10"`
	Password string `json:"password"  form:"password" validate:"min=5,max=10"`
	Phone    string `json:"phone" form:"phone" validate:"required,e164"`
}

//管理端

// Login 登录
func Login(c *gin.Context) {
	user := new(users)
	err := c.Bind(user)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//	判断用户是否存在
	//  查询盐值
	salt := global.Global.Redis.HGet(global.Global.Ctx, user.Username, global.Salt).Val()
	if salt == "" {
		result.Fail(c, global.BadRequest, global.UserNotExistError)
		return
	}
	salts, _ := base64.URLEncoding.DecodeString(salt)
	val := global.Global.Redis.HGet(global.Global.Ctx, user.Username, utils.HashPassword(user.Password, salts)).Val()
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
	use, err := dao.GetUserByUsePwd(user.Username, utils.HashPassword(user.Password, salts))
	if err != nil || use == nil {
		result.Fail(c, global.BadRequest, global.UserNotExistError)
		return
	}
	//生成token
	token := utils.GetToken(use.Identity)
	go func() {
		global.Global.Redis.Set(global.Global.Ctx, use.Identity, token, config.Config.Jwt.Time)
		global.Global.Redis.HSet(global.Global.Ctx, user.Username, utils.HashPassword(user.Password, salts), use.Identity)
	}()
	result.Ok(c, map[string]any{"token": token})
}

// Register 注册
func Register(c *gin.Context) {
	r := new(registers)
	err := c.Bind(r)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//生成随机盐值
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.QueryError)
		return
	}
	//检查用户名是否存在
	username, err := dao.GetUsername(r.Username)
	fmt.Println(username)
	if username != nil {
		result.Fail(c, global.DataConflict, global.QueryError)
		return
	}
	//检查手机号
	username, err = dao.GetPhone(r.Phone)
	if username != nil {
		result.Fail(c, global.DataConflict, global.PhoneError)
		return
	}
	id := utils.GetUidV5(r.Username)
	//插入
	err = dao.InsertUser(&models.User{
		Username: r.Username,
		Identity: id,
		Account:  0,
		Password: utils.HashPassword(r.Password, salt),
		Phone:    r.Phone,
		IP:       c.RemoteIP(),
		Salt:     base64.URLEncoding.EncodeToString(salt),
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataConflict, global.QueryError)
		return
	}
	go func() { //插入identity
		global.Global.Redis.HSet(global.Global.Ctx, r.Username, utils.HashPassword(r.Password, salt), id)
		//盐值
		global.Global.Redis.HSet(global.Global.Ctx, r.Username, global.Salt, base64.URLEncoding.EncodeToString(salt))
	}()
	result.Ok(c, nil)
}

// Info 用户信息
func Info(c *gin.Context) {
	//	获取identity
	id := c.GetString("identity")
	if id == "" {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//
	val := global.Global.Redis.Get(global.Global.Ctx, global.Info+id).Val()
	if val != "" {
		user := new(models.User)
		err := json.Unmarshal([]byte(val), user)
		if err != nil {
			result.Fail(c, global.ServerError, global.ParseError)
			return
		}
		result.Ok(c, user)
		return
	}
	userInfo, err := dao.GetInfoByIdentity(id)
	if err != nil || userInfo == nil {
		result.Fail(c, global.DataNotFound, global.UserNotExistError)
		return
	}
	//同步至redis
	go func() {
		marshal, err := json.Marshal(userInfo)
		if err != nil {
			global.Global.Log.Error(err)
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Info+id, marshal, global.InfoTime*time.Minute).Result()
		if err != nil {
			global.Global.Log.Warn(err)
			return
		}

	}()
	result.Ok(c, userInfo)
}

// Logout 退出登录
func Logout(c *gin.Context) {
	id := c.GetString("identity")
	if id == "" {
		result.Fail(c, global.BadRequest, global.QueryNotFoundError)
		return
	}
	//删除token
	val := global.Global.Redis.Del(global.Global.Ctx, id).Val()
	if val != 1 {
		result.Fail(c, global.ServerError, "退出失败！")
		return
	}
	result.Ok(c, nil)
}

// AssignedAccount 管理员分配账号
func AssignedAccount(c *gin.Context) {
	//输入工号
	uid := c.Query("uid")
	if uid == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//查询uid
	userInfo, err := dao.GetEmployerByUid(uid)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	//生成账号，密码
	//生成随机盐值
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.QueryError)
		return
	}
	//放入redis
	_, err = global.Global.Redis.HSet(global.Global.Ctx, strconv.FormatInt(userInfo.Uid, 10), global.Salt, base64.URLEncoding.EncodeToString(salt)).Result()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.QueryError)
		return
	}

	err = dao.InsertUser(&models.User{
		Username: userInfo.Name,
		Identity: userInfo.Identity,
		Account:  userInfo.Uid,
		Password: utils.HashPassword("123456", salt),
		Phone:    userInfo.Phone,
		IP:       c.RemoteIP(),
		Salt:     base64.URLEncoding.EncodeToString(salt),
	})
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.AssignedAccountError)
		return
	}

	result.Ok(c, map[string]any{
		"account":  userInfo.Uid,
		"password": "123456",
	})

}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	uid := c.Query("id")
	if uid == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//	获取盐值
	val := global.Global.Redis.HGet(global.Global.Ctx, global.Uid, global.Salt).Val()
	if val != "" {
		salt, _ := base64.URLEncoding.DecodeString(val)
		err := dao.UpdatePwd(uid, utils.HashPassword("123456", salt))
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ResetPwdError)
			return
		}
		result.Ok(c, nil)
	}
	//盐值不存在
	account, err := dao.GetByAccount(uid)
	if err != nil || account == nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.BadRequest, global.UserNotFound)
		return
	}
	salt, _ := base64.URLEncoding.DecodeString(account.Salt)
	err = dao.UpdatePwd(uid, utils.HashPassword("123456", salt))
	if err != nil {
		global.Global.Log.Warn(err)
		result.Fail(c, global.BadRequest, global.ResetPwdError)
		return
	}
	result.Ok(c, nil)
}
