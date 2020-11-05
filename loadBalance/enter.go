package main

import (
	"go-gateway/loadBalance/realServer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// start realserver
	startServer("127.0.0.1:2003")
	startServer("127.0.0.1:2004")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func startServer(addr string) {
	rs := &realServer.RealServer{
		Addr: addr,
	}
	rs.Run()
}

//
