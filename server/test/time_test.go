package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	fmt.Println(time.Now().Sub(t1))
	fmt.Println(time.Since(t1))
}
