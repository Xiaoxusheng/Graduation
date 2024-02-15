package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"server/config"
	"server/global"
	"server/result"
	"server/utils"
	"strings"
)

type MyCustomClaims struct {
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

// ParseToken 解析token的中间件
func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	获取唯一标识
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			result.Fail(c, global.AuthenticationFailed, "token 不能为空")
		}
		t := strings.Split(tokenString, " ")

		claims := new(utils.MyCustomClaims)
		token, err := jwt.ParseWithClaims(t[len(t)-1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.Jwt.Key), nil
		})
		_, errs := global.Global.Redis.Get(global.Global.Ctx, claims.Identity).Result()
		if errs != nil {
			fmt.Println(errs)
			result.Fail(c, global.AuthenticationFailed, "token过期或退出登录")
		}

		c.Set("identity", claims.Identity)
		fmt.Println("id", claims.Identity)

		switch {
		case token.Valid:
			fmt.Println("You look nice today")
			c.Next()
		case errors.Is(err, jwt.ErrTokenMalformed):
			fmt.Println("That's not even a token")
			result.Fail(c, global.AuthenticationFailed, "token 错误")
			c.Abort()
			return
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// Invalid signature
			fmt.Println("Invalid signature")
			result.Fail(c, global.AuthenticationFailed, "Invalid signature")
			c.Abort()
			return
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			result.Fail(c, global.AuthenticationFailed, "Timing is everything")
			c.Abort()
			return
		default:
			fmt.Println("Couldn't handle this token:", err)
			result.Fail(c, global.AuthenticationFailed, "Couldn't handle this token")
			c.Abort()
			return
		}
	}
}
