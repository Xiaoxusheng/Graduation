package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"server/config"
	"time"
)

type MyCustomClaims struct {
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

// GetToken 生成token
func GetToken(identity string) string {
	var mySigningKey = []byte(config.Config.Jwt.Key)
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		identity,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.Jwt.Time * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(ss, err)
	return ss
}
