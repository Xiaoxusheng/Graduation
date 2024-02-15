package utils

import (
	"math/rand"
	"time"
)

func GetRandom(m int64) int64 {
	rand.NewSource(time.Now().UnixNano())
	n := rand.Int63n(m)
	if n == 0 {
		n = rand.Int63n(24)
	}
	return n
}

func GetLetter() string {
	rand.NewSource(time.Now().UnixNano())
	s := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n := rand.Int31n(26) + 65
		s[i] = byte(n)
	}
	return string(s)

}

func GetAccount() int64 {
	rand.NewSource(time.Now().UnixNano())
	return rand.Int63n(1000000000)
}
