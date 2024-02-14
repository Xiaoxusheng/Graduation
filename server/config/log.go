package config

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"server/global"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var o = new(sync.Once)

type Log struct {
	io.Writer
	m int64
}

func (l *Log) Format(f *log.Entry) ([]byte, error) {
	var leave int
	switch f.Level {
	case log.InfoLevel, log.DebugLevel:
		leave = global.Gray
	case log.WarnLevel:
		leave = global.Yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		leave = global.Red
	default:
		leave = global.Blue
	}
	var b *bytes.Buffer
	if f.Buffer != nil {
		b = f.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	time1 := f.Time.Format("2006-01-02 15:04:05")
	if f.HasCaller() {
		funcpVal := f.Caller.Function
		fileval := fmt.Sprintf("%s：%d", path.Base(f.Caller.Function), f.Caller.Line)
		fmt.Fprintf(b, "[%s] [%s] [%dkb] \x1b[%dm[%s]\x1b[0m %s %s %s\n", Config.Logs.Prefix, time1, l.m/1024, leave, f.Level, fileval, funcpVal, f.Message)
	} else {
		fmt.Fprintf(b, "[%s] [%s] [%dkb] \x1b[%dm[%s]\x1b[0m %s\n", Config.Logs.Prefix, time1, l.m/1024, leave, leave, f.Level, f.Message)
	}
	return b.Bytes(), nil
}

// 这个是重写
func (l *Log) Write(p []byte) (n int, err error) {
	n, err = l.Writer.Write(p)
	atomic.AddInt64(&l.m, int64(n))
	//fmt.Printf("日志大小 %d kb", l.m/1024)
	return n, err
}

func (l *Log) logFile(m *log.Logger) {
	t := time.Now().Format(time.DateOnly)
	//创建文件
	file, err := os.OpenFile(Config.Logs.Path+t+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	stat, err := file.Stat()
	size := stat.Size()
	l.m = atomic.LoadInt64(&size)
	l.Writer = file
	fmt.Println("create log fail success!")
	s := time.NewTicker(time.Minute * 60 * 24)
	//输出到控制台
	gin.DefaultWriter = io.MultiWriter(l, os.Stdout)
	//日志输出到文件中
	m.SetOutput(io.MultiWriter(l, os.Stdout))
	go func(l *Log, file *os.File) {
		for {
			select {
			case <-s.C:
				//	判读是否超过100m
				if l.m > 100*(1024*1024) {
					err = file.Close()
					if err != nil {
						m.Warn(err)
						return
					}
					t = strings.ReplaceAll(time.Now().Format(time.DateOnly+"-"+time.TimeOnly), ":", "-")
					file, err = os.OpenFile(Config.Logs.Path+t+".log", os.O_CREATE|os.O_RDWR, os.ModePerm)
					if err != nil {
						log.Println(err)
						return
					}
					l = &Log{
						Writer: file,
						m:      0,
					}
					gin.DefaultWriter = io.MultiWriter(l, os.Stdout)
					//输出到控制台,日志文件中
				} else {

					global.Global.Log.Info("进入删除")
					//删除一个月之前日志
					dir, err := os.ReadDir(Config.Logs.Path)
					if err != nil {
						break
					}
					global.Global.Log.Info("dir", dir, len(dir))
					//
					for _, res := range dir {
						t2, err := time.Parse(time.DateOnly, strings.Split(res.Name(), ".")[0])
						if err != nil {
							global.Global.Log.Warn(err)
							continue
						}
						t3 := time.Now().Add(-time.Hour * 24 * 30)
						global.Global.Log.Info(t3.After(t2), t2)
						if t3.After(t2) {
							err := os.Remove(Config.Logs.Path + res.Name())
							if err != nil {
								global.Global.Log.Warn(err)
								continue
							}
							global.Global.Log.Info(res.Name() + "删除")
						}
					}
				}
			}
		}
	}(l, file)
}

func InitLog() {
	o.Do(func() {
		m := log.New()
		l := &Log{
			m: 0,
		}
		////自定义输出
		m.SetFormatter(l)
		//写入文件
		l.logFile(m)
		//输出任务和行号
		m.SetReportCaller(true)
		//最低输出级别
		m.SetLevel(log.InfoLevel)
		global.Global.Log = m
	})
}
