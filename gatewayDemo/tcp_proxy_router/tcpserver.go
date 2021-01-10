package tcp_proxy_router

import (
	"context"
	"gatewayDemo/dao"
	"log"
	"net"
	"strconv"

	"gatewayDemo/reverse_proxy/proxy"
	"gatewayDemo/tcp_proxy_middleware"
	"gatewayDemo/tcp_server"
)

var (
	tcpServerList = []*tcp_server.TcpServer{}
)

type tcpHandler struct{}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TCPServerRun() {

	TCPServiceList := dao.ServiceManagerHandler.GetTCPServiceList()
	for _, Item := range TCPServiceList {
		tempItem := Item
		log.Printf(" [INFO] tcp_proxy_run:%d\n", tempItem.TCPRule.Port)
		// 开启TCP服务器
		go func(serviceDetail *dao.ServiceDetail) {
			addr := ":" + strconv.Itoa(serviceDetail.TCPRule.Port)

			// 获取对应的负载均衡设置
			rb, err := dao.LoadBalanceHandler.GetLoadBalance(serviceDetail)
			if err != nil {
				log.Fatalf(" [ERROR] Get tcp loadbalance err :%s err:%v\n", addr, err)
			}

			//构建路由及设置中间件
			// counter, _ := tcp_proxy_middleware.NewFlowCountService("local_app", time.Second)
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				// 统计层
				tcp_proxy_middleware.TCPFLowCountModeMiddleware(),
				tcp_proxy_middleware.TCPFLowLimitModeMiddleware(),
				// 校验层
				tcp_proxy_middleware.TCPWhiteListModeMiddleware(),
				tcp_proxy_middleware.TCPBlackListModeMiddleware(),
				// tcp_proxy_middleware.FlowCountMiddleWare(counter),
			)

			//构建回调handler
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf(" [ERROR] tcp_proxy_run:%s err:%v\n", addr, err)
			}
		}(tempItem)
	}

}

func TCPServerStop() {
	for _, list := range tcpServerList {
		list.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped \n", list.Addr)
	}
}
