package main

import (
	"fmt"
	"go-gateway/zookeeper/zook"
	"time"
)

func main() {
	zkManager := zook.NewZkManager([]string{"127.0.0.1:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()
	i := 0
	for {
		zkManager.RegistServerPath("/real_server", fmt.Sprint(i))
		fmt.Println("Register", i)
		time.Sleep(5 * time.Second)
		i++
	}
}
