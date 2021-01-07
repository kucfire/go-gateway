package proxy

import (
	"context"
	"fmt"
	"go-gateway/TCP_proxy/server"
	"log"

	"go-gateway/TCP_proxy/load_balance/factory"
	"go-gateway/TCP_proxy/middleware"
	"go-gateway/TCP_proxy/reverse_proxy"
	"net"
)

var (
	addr       = ":2002"
	serveraddr = ":7002"
)

type tcpHandler struct{}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func RunServer() {
	// tcp 服务器测试
	log.Println("Starting tcpserver at : " + serveraddr)
	tcpServe := server.TcpServer{
		Addr:    serveraddr,
		Handler: &tcpHandler{},
	}
	fmt.Println("Starting tcp_server at : " + serveraddr)
	tcpServe.ListenAndServe()
}

// 代理测试
func RunProxy() {
	rb := factory.LoadBalanceFactory(factory.LbWeightRoundRobin)
	rb.Add("127.0.0.1:7002", "40")
	proxy := reverse_proxy.NewTCPLoadBalanceReverseProxy(&middleware.TCPSliceRouterContext{}, rb)

	tcpServe := server.TcpServer{
		Addr:    addr,
		Handler: proxy,
	}
	fmt.Println("starting tcp_proxy at " + addr)
	tcpServe.ListenAndServe()
}

// redis服务器测试
func RunRedisServer() {
	rb := factory.LoadBalanceFactory(factory.LbWeightRoundRobin)
	rb.Add("127.0.0.1:6379", "40")
	proxy := reverse_proxy.NewTCPLoadBalanceReverseProxy(&middleware.TCPSliceRouterContext{}, rb)

	tcpServe := server.TcpServer{
		Addr:    addr,
		Handler: proxy,
	}
	fmt.Println("starting tcp_proxy at " + addr)
	tcpServe.ListenAndServe()
}
