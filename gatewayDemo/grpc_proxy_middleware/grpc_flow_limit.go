package grpc_proxy_middleware

import (
	"errors"
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// 匹配接入方式 给予请求信息
func GRPCFLowLimitModeMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 获取来源的ClientIP
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not find with context")
		}
		split := strings.Split(peerCtx.Addr.String(), ":")
		// fmt.Println(split[0])
		ClientIP := ""
		if len(split) == 2 {
			ClientIP = split[0]
		}

		// fmt.Println(*serviceDetail.AccessControl)
		// 零值校验
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			// 服务端限流
			serviceLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				return err
			}
			if !serviceLimit.Allow() {
				return errors.New(fmt.Sprintf(
					"service flow limit %v",
					serviceDetail.AccessControl.ServiceFlowLimit))
			}

		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			// 客户端限流
			clientLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+ClientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}
			if !clientLimit.Allow() {
				return errors.New(
					fmt.Sprintf("client ip:%v flow limit %v",
						ClientIP, serviceDetail.AccessControl.ClientIPFlowLimit))
			}
		}

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
