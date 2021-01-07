package main

import (
	"go-gateway/TCP_proxy/proxy"
)

func main() {
	// go proxy.RunProxy() // 服务器测试

	// 代理测试
	// proxy.RunServer()

	// redis服务器
	proxy.RunRedisServer()
}
