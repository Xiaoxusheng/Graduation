package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).UnixMilli()
	fmt.Println((time.Now().UnixMilli() - t1) / 1000 / 60 / 60)

}
