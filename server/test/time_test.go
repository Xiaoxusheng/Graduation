package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			fmt.Println(123)
		}
	}

}
