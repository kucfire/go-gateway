package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8074", hystrixStreamHandler)

	// 向circuitSettings写入
	hystrix.ConfigureCommand("kucfireText", // 第一个参数表示的是熔断器的命名
		hystrix.CommandConfig{
			Timeout:                1000, // 单次请求的超时时间，单位mm
			MaxConcurrentRequests:  1,    // 最大并发量
			SleepWindow:            5000, // 熔断后多久去尝试服务是否可用，开启状态下尝试服务是否可转为半打开状态的一个间隔窗口，单位mm
			RequestVolumeThreshold: 10,   // 验证熔断的请求数量，10秒内采样
			ErrorPercentThreshold:  10,   // 验证熔断的错误百分比，当达到百分百的时候，则进入熔断开启状态
		})

	for i := 0; i < 10000; i++ {
		err := hystrix.Do("kucfireText", func() error {
			// test case 1 并发测试
			if i == 0 {
				return errors.New("server error")
			}
			// test case 2 超时测试
			// time.Sleep(2 * time.Second)
			log.Println("do server")
			return nil
		}, func(err error) error {
			if err != nil {
				if err == errors.New("server error") {
					return nil
				} else if err == errors.New("max concurrency") {
					return nil
				}
				return errors.New("other error")
			}
			return nil
		}) // feedback 可以为nil，一般是捕获hystrix.runFunc的错误进行处理
		if err != nil {
			log.Println("hystrix err:" + err.Error())
			time.Sleep(1 * time.Second)
			log.Println("Sleep 1 Second")
		}
	}

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
