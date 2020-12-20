package main

import (
	"go-gateway/httpsAndhttp2/http2proxy"
	"go-gateway/httpsAndhttp2/realserver"
	"os"
	"os/signal"
	"syscall"
)

func runRealServer() {
	realServer1 := &realserver.RealServer{Addr: "127.0.0.1:3003"}
	realServer1.Run()
	realServer2 := &realserver.RealServer{Addr: "127.0.0.1:3004"}
	realServer2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func runProxy() {
	// httpsProxy.Run()
	http2proxy.Run()
}

func main() {
	go runRealServer()
	go runProxy()
	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
