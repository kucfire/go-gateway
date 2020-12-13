package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// 定义一个新的限制器
	l := rate.NewLimiter(1, 5)
	fmt.Println(l.Limit(), l.Burst())

	//
	for i := 0; i < 100; i++ {
		// 阻塞，等待直到去到一个token
		log.Println("before wait")
		c, _ := context.WithTimeout(context.Background(), time.Second*2)
		if err := l.Wait(c); err != nil {
			log.Println("limiter wait err:" + err.Error())
		}
		log.Println("after wait")

		// 返回需要等待多久才有新的token，这样可以等待指定时间执行任务
		r := l.Reserve()
		log.Println("reverse Delay", r.Delay())

		// 判断当前是否可以取到token
		a := l.Allow()
		log.Println("Allow:", a)

		// time.Sleep(200 * time.Millisecond)
		// log.Println(time.Now().Format("2006-01-02 15:04:05.000"))
	}
}
