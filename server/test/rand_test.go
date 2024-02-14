package test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// HashPassword 使用SHA256和盐加密密码
func HashPassword(password string, salt []byte) string {
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(password))
	sha256Hasher.Write(salt)
	hashedPassword := sha256Hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedPassword)
}

func Test_rand(t *testing.T) {
	salt, err := GenerateSalt(16)
	if err != nil {
		return
	}
	d := HashPassword("njdjs", salt)
	fmt.Println(d)
}
