package test

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			client := http.Client{
				Timeout: 10 * time.Second,
			}
			get, err := http.NewRequest(http.MethodGet, "http://localhost:80/user/get_notice_list/", nil)
			if err != nil {
				t.Error(err.Error())
			}
			get.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6IjBhZDM0Zjk5LWM1NGQtNDgwOC04NTI3LTUzMTJlZTlhZmIzOCIsImV4cCI6MTcxNDQ0OTcyOSwiaWF0IjoxNzE0MzYzMzI5fQ.UYizXkcK_r6UkLQSbJYgf5T-HFPmjoLv5QXuVjEO1UQ")
			resp, err := client.Do(get)
			if err != nil {
				t.Error(err.Error())
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err.Error())
			}
			if resp.StatusCode != 200 {
				t.Error("请求失败", resp.StatusCode, string(body))
			}
			t.Log(string(body))
		}()
	}
	select {}

}
