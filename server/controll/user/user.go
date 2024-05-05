package user

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required" `
}

// Login /*登录要保证逻辑上的并发安全*/
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
	//新密码
	pwd := c.Query("password")
	//旧密码
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
		//判断旧密码是否正确
		//判断hash中是否有这条记录
		if !global.Global.Redis.HExists(global.Global.Ctx, uid, utils.HashPassword(p, salt)).Val() {
			//	旧密码错误
			result.Fail(c, global.BadRequest, global.OldPedError)
			return
		}
		//更新
		err := dao.UpdatePwd(uid, utils.HashPassword(pwd, salt))
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ChangePwdError)
			return
		}
		//删除旧的
		_, err = global.Global.Redis.HDel(global.Global.Ctx, uid, utils.HashPassword(p, salt)).Result()
		if err != nil {
			global.Global.Log.Warn(err)
			result.Fail(c, global.BadRequest, global.ChangePwdError)
			return
		}
		//设置新的
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
	//判断旧密码是否正确
	if !global.Global.Redis.HExists(global.Global.Ctx, uid, utils.HashPassword(p, salt)).Val() {
		//	旧密码错误
		result.Fail(c, global.BadRequest, global.OldPedError)
		return
	}
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

// EmployeeInfo 员工个人人信息
func EmployeeInfo(c *gin.Context) {
	id := c.GetString("identity")
	if id == "" {
		global.Global.Log.Warn("identity is null")
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	var uid string
	uid = global.Global.Redis.HGet(global.Global.Ctx, global.UidId, id).Val()
	if uid == "" {
		employer, err := dao.GetUserById(id)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.UserNotExistError)
			return
		}
		//判断员工是否在职
		if employer.Status != 1 {
			result.Fail(c, global.DataConflict, global.UserNotWorkError)
			return
		}
		uid = strconv.FormatInt(employer.Uid, 10)
		//	同步到redis
		err = global.Global.Pool.Submit(func() {
			global.Global.Wg.Add(1)
			defer global.Global.Wg.Done()
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.UidId, id, employer.Uid).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
		})
		if err != nil {
			global.Global.Log.Error("submit err :", err)
			result.Fail(c, global.ServerError, global.GetClockError)
			return
		}
		return
	}
	//获取员工信息
	val := global.Global.Redis.Get(global.Global.Ctx, global.Uid+uid).Val()
	if val != "" {
		e := new(models.Employee)
		err := json.Unmarshal([]byte(val), e)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		result.Ok(c, e)
		return
	}
	//不存在，过期
	atoi, err := strconv.Atoi(uid)
	if err != nil {
		result.Fail(c, global.DataUnmarshal, global.AtoiError)
		return
	}
	info, err := dao.GetEmployerInfo(int64(atoi))
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.DataNotFound, global.UserNotExistError)
		return
	}
	//插入
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		if info == nil {
			_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+uid, "null", global.InfoTime*time.Second).Result()
			return
		}
		marshal, err := json.Marshal(info)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Uid+uid, marshal, global.InfoTime*time.Second).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	result.Ok(c, info)
}

// ChangeInfo 修改个人信息
func ChangeInfo(c *gin.Context) {
	e := new(global.Infos)
	err := c.Bind(e)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryError)
		global.Global.Log.Error(err)
		return
	}
	if e.Uid == 0 || e.Sex == 0 {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//判断员工是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.Employer, e.Uid).Val()
	if !val {
		result.Fail(c, global.BadRequest, global.EmployerNotFoundError)
		return
	}
	err = global.Global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//
		err = dao.UpdateUserinfo(tx, e)
		if err != nil {
			result.Fail(c, global.ServerError, global.UpdateUserInfoError)
			global.Global.Log.Error(err)
			return err
		}
		err = dao.UpdateUser(tx, strconv.FormatInt(e.Uid, 10), e.Name)
		if err != nil {
			result.Fail(c, global.ServerError, global.UpdateUserInfoError)
			global.Global.Log.Error(err)
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	//	更新信息
	if err != nil {
		result.Fail(c, global.ServerError, global.UpdateUserInfoError)
		global.Global.Log.Error(err)
		return
	}
	//删除缓存中信息
	err = global.Global.Pool.Submit(func() {
		global.Global.Wg.Add(1)
		defer global.Global.Wg.Done()
		_, err := global.Global.Redis.Del(global.Global.Ctx, global.Uid+strconv.Itoa(int(e.Uid))).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	})
	result.Ok(c, nil)
}
