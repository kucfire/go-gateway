package tcpserver

import (
	"context"
	"fmt"
	"gatewayDemo/tcp_server"
	"log"
	"net"
)

var (
	addr = ":7002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TCPServerRun() {
	//tcp服务器测试
	log.Println("Starting tcpserver at " + addr)
	tcpServ := &tcp_server.TcpServer{
		Addr:    addr,
		Handler: &tcpHandler{},
	}
	fmt.Println("Starting tcp_server at " + addr)
	tcpServ.ListenAndServe()
}
