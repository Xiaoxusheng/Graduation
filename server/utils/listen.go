package utils

import (
	"os"
	"os/signal"
	"server/global"
	"syscall"
	"time"
)

func Listen() {

	signalCh := make(chan os.Signal, 1)
	// 监听指定的系统信号
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 来处理接收到的系统信号
	// 接收信号
	for {
		select {
		case sig := <-signalCh:
			global.Global.Log.Printf("Received signal: %s\n", sig)
			//关闭连接
			db, err := global.Global.Mysql.DB()
			if err != nil {
				return
			}
			global.Global.Log.Info("mysql connect close success")
			err = db.Close()
			if err != nil {
				global.Global.Log.Warn("mysql close err:", err)
			}
			err = global.Global.Redis.Close()
			if err != nil {
				global.Global.Log.Warn("redis close err:", err)
			}
			global.Global.Log.Info("redis connect close success")
			err = global.Global.Pool.ReleaseTimeout(time.Second * 5)
			if err != nil {
				global.Global.Log.Error(err)
				return
			}
			global.Global.Log.Info("pool close success")

			//退出
			os.Exit(0)
		}

	}
	// 打印接收到的信号
	// 进行一些清理工作或其他操作
	// ...

	// 退出程序
	//os.Exit(0)
}
