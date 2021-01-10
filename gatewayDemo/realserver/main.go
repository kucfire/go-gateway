package main

import (
	"flag"
	"gatewayDemo/realserver/tcpserver"

	"gatewayDemo/realserver/httpserver"
	"os"
)

var (
	realserver = flag.Int(
		"servercode", // flag name 标志命名
		0,            // flag value 值
		"input 1(http server) or 2(tcp server) or 3(grpc server)", // flag usge 用法
	)
)

func main() {
	flag.Parse()
	if *realserver == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *realserver == 1 {
		httpserver.HTTPServerRun()
	}

	if *realserver == 2 {
		tcpserver.TCPServerRun()
	}
	if *realserver == 3 {

	}
}
