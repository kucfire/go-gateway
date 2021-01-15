package proxy

import (
	"context"
	"fmt"

	// "go-gateway/GRPC_proxy/reverse_proxy_lb/load_balance/config"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"log"

	// "go-gateway/GRPC_proxy/grpc-proxy/proxy"
	"gatewayDemo/grpc_proxy/proxy"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGrpcLoadBalanceHandler(lb config.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		fmt.Println("Addr : ", nextAddr)
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}

		// // 重写上下文
		// director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		// 	// 拒绝某些特殊请求
		// 	if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
		// 		return ctx, nil, status.Errorf(codes.Unimplemented,
		// 			"Unknown method")
		// 	}
		// 	c, err := grpc.DialContext(ctx, "localhost:50055", grpc.WithCodec(grpc_proxy.Codec()), grpc.WithInsecure())
		// 	md, _ := metadata.FromIncomingContext(ctx)
		// 	outCtx, _ := context.WithCancel(ctx)
		// 	outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		// 	return outCtx, c, err
		// }
		return proxy.TransparentHandler(director)
	}()
}
