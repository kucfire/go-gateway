package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	gw "gatewayDemo/realserver/grpcserver/GRPC/proto"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	serverAddr         = ":8081"
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "127.0.0.1:50055", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 可以理解为每个rs都需要持续跟下游建立连接
	mux := runtime.NewServeMux() //设置路由
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterEchoHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(serverAddr, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	fmt.Println("server listening at", serverAddr)
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
