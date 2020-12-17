package main

import (
	"go-gateway/websocket/proxy"
	"go-gateway/websocket/testserver"
	"os"
	"os/signal"
	"syscall"

)

func runRealServer() {
	realServer1 := &testserver.RealServer{Addr: "127.0.0.1:2003"}
	realServer1.Run()
	realServer2 := &testserver.RealServer{Addr: "127.0.0.1:2004"}
	realServer2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func runProxy() {
	proxy.Run()
}

func main() {
	go runProxy()
	runRealServer()
}
