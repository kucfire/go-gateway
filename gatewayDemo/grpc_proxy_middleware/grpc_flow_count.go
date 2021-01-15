package grpc_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/public"

	"google.golang.org/grpc"
)

// 匹配接入方式 给予请求信息
func GRPCFLowCountModeMiddleware(serviceDetail *dao.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 全站统计
		totalCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowTotalPrefix)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		// 服务统计
		serviceCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := handler(srv, ss); err != nil {
			// log.Printf("RPC failed with error %v\n", err)
			return err
		}

		return nil
	}
}
