package grpc_proxy_router

import (
	"gatewayDemo/dao"
	grpc_proxy "gatewayDemo/grpc_proxy/proxy"

	// "gatewayDemo/grpc_proxy_middleware"
	"gatewayDemo/grpc_proxy_middleware"
	"gatewayDemo/reverse_proxy/proxy"

	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

var (
	grpcServerList = []*warpGRPCServer{}
)

type warpGRPCServer struct {
	Addr string
	*grpc.Server
}

func GRPCServerRun() {
	GRPCServiceList := dao.ServiceManagerHandler.GetGRPCServiceList()
	for _, Item := range GRPCServiceList {
		tempItem := Item
		// 开启GRPC服务器
		go func(serviceDetail *dao.ServiceDetail) {
			addr := ":" + strconv.Itoa(serviceDetail.GRPCRule.Port)

			// 监听tcp端口
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [ERROR] GRPC failed to listen : %s err : %v\n", addr, err)
			}

			// 获取对应的负载均衡设置
			rb, err := dao.LoadBalanceHandler.GetLoadBalance(serviceDetail)
			if err != nil {
				log.Fatalf(" [ERROR] Get GRPC loadbalance err : %s err : %v\n", addr, err)
			}

			grpcHandler := proxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor( // 中间件层
					grpc_proxy_middleware.GRPCAuthStreamInterceptor,
					grpc_proxy_middleware.GRPCFLowCountModeMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCFLowLimitModeMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJwtOauthTokenModeMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCJwtFLowCountModeMiddleware(),
					grpc_proxy_middleware.GRPCFJwtLowLimitModeMiddleware(),
					grpc_proxy_middleware.GRPCWhiteListModeMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCBlackListModeMiddleware(serviceDetail),
					grpc_proxy_middleware.GRPCHeaderTransferModeMiddleware(serviceDetail)),
				grpc.CustomCodec(grpc_proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)

			grpcServerList = append(grpcServerList, &warpGRPCServer{Addr: addr, Server: s})
			log.Printf("[INFO] grpc_proxy_run %v\n", addr)
			if err := s.Serve(lis); err != nil {
				log.Fatalf(" [ERROR] grpc_proxy_run : %s err : %v\n", addr, err)
			}
		}(tempItem)
	}
}

func GRPCServerStop() {
	for _, list := range grpcServerList {
		list.GracefulStop()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped \n", list.Addr)
	}
}
