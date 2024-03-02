package test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"
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
	s := base64.URLEncoding.EncodeToString(salt)
	a, _ := base64.URLEncoding.DecodeString(s)
	fmt.Println(salt, s, a, d)

	// 获取当前时间
	now := time.Now()

	// 获取当前月份的天数
	//year, month, _ := now.Date()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	fmt.Println("Days in current month:", daysInMonth, daysInMonth.Day())

	// 要转换的日期字符串
	dateStr := "2024-06"

	// 定义日期格式
	layout := "2006-01"

	// 将字符串转换为时间类型
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// 打印转换后的时间
	fmt.Println("Converted date:", date)

	ss := now.Add(time.Hour * 5)
	fmt.Println(ss.Sub(now).Hours())
}
