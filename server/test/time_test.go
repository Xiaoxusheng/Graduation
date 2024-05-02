package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	a := [3]int{0, 1, 2}
	Change(a)
	fmt.Println(a)
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			fmt.Println(123)
		}
	}

}

func Change(a [3]int) {
	for i, v := range a {
		a[2] = 10
		if i == 2 {
			fmt.Println(a)
			fmt.Println(v)
		}
	}
}
